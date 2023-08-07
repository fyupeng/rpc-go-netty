package codec

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty/codec"
	"github.com/go-netty/go-netty/utils"
	"log"
	"rpc-go-netty/protocol"
	"rpc-go-netty/serializer"
)

const (
	MagicNumberLength         = 2
	PackageProtocolTypeLength = 1
	SerializationTypeLength   = 1
	DataLengthTypeLength      = 4
)

func CommonCodec(lengthFieldOffset int, lengthFieldLength int, serializationType int) codec.Codec {
	utils.AssertIf(lengthFieldOffset < 0, "maxFrameLength must be a positive integer")
	utils.AssertIf(lengthFieldLength <= 0, "delimiter must be nonempty string")
	return &commonCodec{
		magicNumber:       0xBABE,
		serializationType: serializationType,
		lengthFieldOffset: lengthFieldOffset,
		lengthFieldLength: lengthFieldLength,
	}
}

type commonCodec struct {
	magicNumber       int // 魔数
	serializationType int // 序列化类型
	lengthFieldOffset int // 协议头偏移
	lengthFieldLength int // 协议头长度
}

func (*commonCodec) CodecName() string {
	return "common-codec"
}

func (codec *commonCodec) HandleWrite(ctx netty.OutboundContext, message netty.Message) {

	fmt.Println("codec 准备写入数据，正在编码")

	fmt.Println("message ", message)

	fmt.Println(message)

	// 构建协议头字节流
	protocolHeader := make([]byte, codec.lengthFieldLength)
	// 设置魔数（Magic Number）
	buffIdx := codec.lengthFieldOffset
	binary.BigEndian.PutUint16(protocolHeader[buffIdx:buffIdx+MagicNumberLength], uint16(codec.magicNumber))
	buffIdx += MagicNumberLength

	// 设置协议包类型
	packageProtocolType := transProtocolCode(message)

	switch PackageProtocolTypeLength {
	case 1:
		protocolHeader[buffIdx] = byte(packageProtocolType)
	case 2:
		binary.BigEndian.PutUint16(protocolHeader[buffIdx:buffIdx+PackageProtocolTypeLength], uint16(packageProtocolType))
	case 3:
		binary.BigEndian.PutUint32(protocolHeader[buffIdx:buffIdx+PackageProtocolTypeLength], uint32(packageProtocolType))
	case 4:
		binary.BigEndian.PutUint64(protocolHeader[buffIdx:buffIdx+PackageProtocolTypeLength], uint64(packageProtocolType))
	default:
		binary.BigEndian.PutUint16(protocolHeader[buffIdx:buffIdx+PackageProtocolTypeLength], uint16(packageProtocolType))
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
	utils.AssertIf(err != nil, "serialize fatail: ", err)

	dataLength := len(serializedDataBytes)
	fmt.Println("序列化后的数据长度：", dataLength)
	binary.BigEndian.PutUint32(protocolHeader[buffIdx:buffIdx+DataLengthTypeLength], uint32(dataLength))

	fmt.Println(protocolHeader)

	// 将协议头和数据部分拼接，形成最终的编码后的字节流
	encodedData := append(protocolHeader, serializedDataBytes...)
	//encodedData = append(encodedData, '$')

	ctx.HandleWrite(encodedData)
}

func (codec *commonCodec) HandleRead(ctx netty.InboundContext, message netty.Message) {

	fmt.Println("common_codec 正在读取数据，正在解码")

	dataBytesa := utils.MustToBytes(message)

	fmt.Println(dataBytesa)

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

	packageProtocol := getProtocolByCode(BytesToInt(readBuff))

	fmt.Println("123")
	fmt.Println(packageProtocol)

	fmt.Println(&protocol.RpcResponseProtocol1{})

	utils.AssertIf(packageProtocol == nil, "Invalid package protocol type:%X", readBuff)

	// 读取序列化协议（1字节）
	readBuff = make([]byte, 0, SerializationTypeLength)
	for len(readBuff) < SerializationTypeLength {
		n := utils.AssertLength(reader.Read(tempBuff[:]))
		readBuff = append(readBuff, tempBuff[:n]...)
	}
	serializaer := GetSerializerByCode(codec.serializationType)
	utils.AssertIf(serializaer == nil, "Invalid serializer type:%X", readBuff)
	// 读取数据长度（4字节）
	readBuff = make([]byte, 0, DataLengthTypeLength)
	for len(readBuff) < DataLengthTypeLength {
		n := utils.AssertLength(reader.Read(tempBuff[:]))
		readBuff = append(readBuff, tempBuff[:n]...)
	}
	dataLength := BytesToInt(readBuff)

	fmt.Println("dataLength ", dataLength)

	// 读取数据（${dataLength} 字节）
	readBuff = make([]byte, 0, dataLength)
	for len(readBuff) < dataLength {
		n := utils.AssertLength(reader.Read(tempBuff[:]))
		readBuff = append(readBuff, tempBuff[:n]...)
	}
	dataBytes := readBuff

	// 根据序列化协议 反序列化
	data0, err0 := serializaer.Deserialize(dataBytes, &protocol.RpcResponseProtocol1{})

	fmt.Println("data0 ", data0)

	if err0 != nil {
		log.Fatal("err0: ", err0)
	}

	data, err1 := serializaer.Deserialize(dataBytes, packageProtocol)

	fmt.Println("data ", data)

	if err1 != nil {
		log.Fatal("dataBytes deserialize failed：", err1)
	}

	// post message
	ctx.HandleRead(packageProtocol)

}

type Result struct {
	Code        int    `json:"code"`
	Msg         string `json:"msg"`
	Description string `json:"description"`
}

func GetSerializerByCode(serializationTypeCode int) (serial serializer.CommonSerializer) {
	switch serializationTypeCode {
	case serializer.KryoSerializerCode:
		serial = serializer.NewKryoSerializer()
	case serializer.JsonSerializerCode:
		serial = serializer.NewJsonSerializer()
	case serializer.HessianSerializerCode:
		serial = serializer.NewHessianSerializer()
	default:
	}
	return
}

func getProtocolByCode(protocolCode int) (proto protocol.Protocol) {
	switch protocolCode {
	case protocol.RequestProtocolCode:
		proto = protocol.NewRpcRequestProtocol()
	case protocol.ResponseProtocolCode:
		proto = protocol.NewRpcResponseProtocol()
	case protocol.UnRecognizeProtocolCode:
		log.Println("unrecognized protocol:", protocolCode)
	default:
		log.Println("unrecognized protocol:", protocolCode)
	}
	return
}

func transProtocolCode(proto protocol.Protocol) (protocolCode int) {
	switch proto.(type) {
	case protocol.RequestProtocol:
		// 处理请求逻辑
		protocolCode = protocol.RequestProtocolCode
	case protocol.ResponseProtocol:
		// 处理响应逻辑
		protocolCode = protocol.ResponseProtocolCode
	default:
		// 处理其他情况
		protocolCode = protocol.UnRecognizeProtocolCode
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
