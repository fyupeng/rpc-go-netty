package codec

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty/codec"
	"github.com/go-netty/go-netty/utils"
	"log"
	"reflect"
	"rpc-go-netty/protocol"
	"rpc-go-netty/serializer"
)

const (
	MagicNumberLength         = 2
	PackageProtocolTypeLength = 1
	SerializationTypeLength   = 1
	DataLengthTypeLength      = 4
)

func CommonCodec(lengthFieldOffset int, lengthFieldLength int, packageProtocolType int, serializationType int) codec.Codec {
	utils.AssertIf(lengthFieldOffset < 0, "maxFrameLength must be a positive integer")
	utils.AssertIf(lengthFieldLength <= 0, "delimiter must be nonempty string")
	return &commonCodec{
		magicNumber:         0xBABE,
		packageProtocolType: packageProtocolType,
		serializationType:   serializationType,
		lengthFieldOffset:   lengthFieldOffset,
		lengthFieldLength:   lengthFieldLength,
	}
}

type commonCodec struct {
	magicNumber         int // 魔数
	packageProtocolType int // 包协议类型
	serializationType   int // 序列化类型
	lengthFieldOffset   int // 协议头偏移
	lengthFieldLength   int // 协议头长度
}

func (*commonCodec) CodecName() string {
	return "common-codec"
}

func (codec *commonCodec) HandleWrite(ctx netty.OutboundContext, message netty.Message) {

	fmt.Println("codec 准备写入数据，正在编码")

	fmt.Println("message ", message)

	// 构建协议头字节流
	protocolHeader := make([]byte, codec.lengthFieldLength)
	// 设置魔数（Magic Number）
	buffIdx := codec.lengthFieldOffset
	binary.BigEndian.PutUint16(protocolHeader[buffIdx:buffIdx+MagicNumberLength], uint16(codec.magicNumber))
	buffIdx += MagicNumberLength
	// 设置协议包类型
	switch PackageProtocolTypeLength {
	case 1:
		protocolHeader[buffIdx] = byte(codec.packageProtocolType)
	case 2:
		binary.BigEndian.PutUint16(protocolHeader[buffIdx:buffIdx+PackageProtocolTypeLength], uint16(codec.packageProtocolType))
	case 3:
		binary.BigEndian.PutUint32(protocolHeader[buffIdx:buffIdx+PackageProtocolTypeLength], uint32(codec.packageProtocolType))
	case 4:
		binary.BigEndian.PutUint64(protocolHeader[buffIdx:buffIdx+PackageProtocolTypeLength], uint64(codec.packageProtocolType))
	default:
		binary.BigEndian.PutUint16(protocolHeader[buffIdx:buffIdx+PackageProtocolTypeLength], uint16(codec.packageProtocolType))
	}
	buffIdx += PackageProtocolTypeLength
	// 设置序列化类型
	switch SerializationTypeLength {
	case 1:
		protocolHeader[buffIdx] = byte(codec.serializationType)
	case 2:
		binary.BigEndian.PutUint16(protocolHeader[buffIdx:buffIdx+SerializationTypeLength], uint16(codec.serializationType))
	case 3:
		binary.BigEndian.PutUint32(protocolHeader[buffIdx:buffIdx+SerializationTypeLength], uint32(codec.serializationType))
	case 4:
		binary.BigEndian.PutUint64(protocolHeader[buffIdx:buffIdx+SerializationTypeLength], uint64(codec.serializationType))
	default:
		binary.BigEndian.PutUint16(protocolHeader[buffIdx:buffIdx+PackageProtocolTypeLength], uint16(codec.serializationType))
	}
	buffIdx += SerializationTypeLength

	// 设置数据长度 - 先序列化数据，再计算长度
	serial := GetSerializerByCode(codec.serializationType)

	serializedDataBytes, err := serial.Serialize(message)

	dataLength := len(serializedDataBytes)
	fmt.Println("序列化后的数据长度：", dataLength)
	binary.BigEndian.PutUint32(protocolHeader[buffIdx:buffIdx+DataLengthTypeLength], uint32(dataLength))

	// 将协议头和数据部分拼接，形成最终的编码后的字节流

	testDeserialize := protocol.NewRpcRequestProtocol()
	err1 := serial.Deserialize(serializedDataBytes, &testDeserialize)
	if err1 != nil {
		log.Fatal("err1: ", err1)
	}

	utils.AssertIf(err != nil, "serialize fatail: ", err)

	encodedData := append(protocolHeader, serializedDataBytes...)

	ctx.HandleWrite(encodedData)
}

func (codec *commonCodec) HandleRead(ctx netty.InboundContext, message netty.Message) {

	fmt.Println("common_codec 正在读取数据，正在解码")

	dataBytesa := utils.MustToBytes(message)

	reader := utils.MustToReader(dataBytesa)

	// 读取魔数（2字节）
	readBuff := make([]byte, 0, MagicNumberLength)
	tempBuff := make([]byte, 1)
	for len(readBuff) < MagicNumberLength {
		n := utils.AssertLength(reader.Read(tempBuff[:]))
		readBuff = append(readBuff, tempBuff[:n]...)
	}
	protocolCheckPass := bytes.Equal(IntToBytes(codec.magicNumber), readBuff)

	utils.AssertIf(!protocolCheckPass, "Invalid magic number:%s", readBuff)
	if !protocolCheckPass {
		return
	}

	// 读取协议头字段（1字节）
	readBuff = make([]byte, 0, PackageProtocolTypeLength)
	for len(readBuff) < PackageProtocolTypeLength {
		n := utils.AssertLength(reader.Read(tempBuff[:]))
		readBuff = append(readBuff, tempBuff[:n]...)
	}
	packageProtocolType := getPackageProtocolTypeByCode(BytesToInt(readBuff))
	utils.AssertIf(packageProtocolType == nil, "Invalid package protocol type:%X", readBuff)

	// 读取序列化协议（1字节）
	readBuff = make([]byte, 0, SerializationTypeLength)
	for len(readBuff) < SerializationTypeLength {
		n := utils.AssertLength(reader.Read(tempBuff[:]))
		readBuff = append(readBuff, tempBuff[:n]...)
	}
	serializationType := GetSerializerByCode(codec.serializationType)
	utils.AssertIf(serializationType == nil, "Invalid serializer type:%X", readBuff)
	// 读取数据长度（4字节）
	readBuff = make([]byte, 0, DataLengthTypeLength)
	for len(readBuff) < DataLengthTypeLength {
		n := utils.AssertLength(reader.Read(tempBuff[:]))
		readBuff = append(readBuff, tempBuff[:n]...)
	}
	dataLength := BytesToInt(readBuff)

	// 读取数据（${dataLength} 字节）
	readBuff = make([]byte, 0, dataLength)
	for len(readBuff) < dataLength {
		n := utils.AssertLength(reader.Read(tempBuff[:]))
		readBuff = append(readBuff, tempBuff[:n]...)
	}
	dataBytes := readBuff

	// 根据序列化协议 反序列化
	deserializedObject := protocol.NewRpcRequestProtocol()
	err := serializationType.Deserialize(dataBytes, &deserializedObject)

	if err != nil {
		log.Fatal("serializationType["+TransSerializationType(serializationType.GetValue())+"] deserialize failed：", err)
	}

	fmt.Println("反序列化成功： ", deserializedObject)

	// post message
	ctx.HandleRead(deserializedObject)

}

type Result struct {
	Code        int    `json:"code"`
	Msg         string `json:"msg"`
	Description string `json:"description"`
}

func GetSerializerByCode(serializationTypeCode int) (serial serializer.CommonSerializer) {
	switch serializationTypeCode {
	case 0:
		serial = serializer.KryoSerializer()
	case 1:
		serial = serializer.JsonSerializer()
	case 2:
		serial = serializer.HessianSerializer()
	default:
	}
	return
}

func TransSerializationType(serializationType int) (serial string) {
	switch serializationType {
	case 0:
		serial = "hryo"
	case 1:
		serial = "json"
	case 2:
		serial = "hessian"
	default:
	}
	return
}

func getPackageProtocolTypeByCode(protocolTypeCode int) (typ reflect.Type) {
	switch protocolTypeCode {
	case 71:
		typ = reflect.TypeOf(protocol.NewRpcRequestProtocol())
	case 73:
		typ = reflect.TypeOf(protocol.NewRpcResponseProtocol())
	default:
	}
	return
}

func transPackageProtocolType(protocolType int) (proto string) {
	switch protocolType {
	case 71:
		proto = "q"
	case 73:
		proto = "s"
	default:
	}
	return
}

// 整形转换成字节
func IntToBytes(val int) []byte {
	byteArray := make([]byte, 0)
	for val > 0 {
		byteArray = append(byteArray, byte(val%256))
		val /= 256
	}
	// 反转字节数组
	for i, j := 0, len(byteArray)-1; i < j; i, j = i+1, j-1 {
		byteArray[i], byteArray[j] = byteArray[j], byteArray[i]
	}
	return byteArray
}

func BytesToInt(bytes []byte) int {
	result := 0
	for _, b := range bytes {
		result = (result << 8) | int(b)
	}
	return result
}
