package codec

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty/codec"
	"github.com/go-netty/go-netty/utils"
	"log"
	"rpc-go-netty/protocol"
	"rpc-go-netty/serializer"
)

const (
	MagicNumberLength       = 2
	ProtocolCodeLength      = 1
	SerializationCodeLength = 1
	DataLengthCodeLength    = 4
)

func CommonCodec(lengthFieldOffset int, lengthFieldLength int, serializerCode int) codec.Codec {
	utils.AssertIf(lengthFieldOffset < 0, "maxFrameLength must be a positive integer")
	utils.AssertIf(lengthFieldLength <= 0, "delimiter must be nonempty string")
	return &commonCodec{
		magicNumber:       0xBABE,
		serializerCode:    serializerCode,
		lengthFieldOffset: lengthFieldOffset,
		lengthFieldLength: lengthFieldLength,
	}
}

type commonCodec struct {
	magicNumber       int // 魔数
	serializerCode    int // 序列化类型
	lengthFieldOffset int // 协议头偏移
	lengthFieldLength int // 协议头长度
}

func (*commonCodec) CodecName() string {
	return "common-codec"
}

func (codec *commonCodec) HandleWrite(ctx netty.OutboundContext, message netty.Message) {

	// 构建协议头字节流
	protocolHeader := make([]byte, codec.lengthFieldLength)
	// 设置魔数（Magic Number）
	buffIdx := codec.lengthFieldOffset
	binary.BigEndian.PutUint16(protocolHeader[buffIdx:buffIdx+MagicNumberLength], uint16(codec.magicNumber))
	buffIdx += MagicNumberLength

	// 设置协议包类型
	protocolCode := transProtocolCode(message)

	switch ProtocolCodeLength {
	case 1:
		protocolHeader[buffIdx] = byte(protocolCode)
	case 2:
		binary.BigEndian.PutUint16(protocolHeader[buffIdx:buffIdx+ProtocolCodeLength], uint16(protocolCode))
	case 3:
		binary.BigEndian.PutUint32(protocolHeader[buffIdx:buffIdx+ProtocolCodeLength], uint32(protocolCode))
	case 4:
		binary.BigEndian.PutUint64(protocolHeader[buffIdx:buffIdx+ProtocolCodeLength], uint64(protocolCode))
	default:
		binary.BigEndian.PutUint16(protocolHeader[buffIdx:buffIdx+ProtocolCodeLength], uint16(protocolCode))
	}
	buffIdx += ProtocolCodeLength
	// 设置序列化类型
	switch ProtocolCodeLength {
	case 1:
		protocolHeader[buffIdx] = byte(codec.serializerCode)
	case 2:
		binary.BigEndian.PutUint16(protocolHeader[buffIdx:buffIdx+SerializationCodeLength], uint16(codec.serializerCode))
	case 3:
		binary.BigEndian.PutUint32(protocolHeader[buffIdx:buffIdx+SerializationCodeLength], uint32(codec.serializerCode))
	case 4:
		binary.BigEndian.PutUint64(protocolHeader[buffIdx:buffIdx+SerializationCodeLength], uint64(codec.serializerCode))
	default:
		binary.BigEndian.PutUint16(protocolHeader[buffIdx:buffIdx+SerializationCodeLength], uint16(codec.serializerCode))
	}
	buffIdx += SerializationCodeLength

	// 设置数据长度 - 先序列化数据，再计算长度
	serial := GetSerializerByCode(codec.serializerCode)

	serializedDataBytes, err := serial.Serialize(message)
	utils.AssertIf(err != nil, "serialize fatail: ", err)

	dataLength := len(serializedDataBytes)
	binary.BigEndian.PutUint32(protocolHeader[buffIdx:buffIdx+DataLengthCodeLength], uint32(dataLength))

	// 将协议头和数据部分拼接，形成最终的编码后的字节流
	encodedData := append(protocolHeader, serializedDataBytes...)
	//encodedData = append(encodedData, '$')

	ctx.HandleWrite(encodedData)
}

func (codec *commonCodec) HandleRead(ctx netty.InboundContext, message netty.Message) {

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
	readBuff = make([]byte, 0, ProtocolCodeLength)
	for len(readBuff) < ProtocolCodeLength {
		n := utils.AssertLength(reader.Read(tempBuff[:]))
		readBuff = append(readBuff, tempBuff[:n]...)
	}

	packageProtocol := getProtocolByCode(BytesToInt(readBuff))

	utils.AssertIf(packageProtocol == nil, "Invalid package protocol type:%X", readBuff)

	// 读取序列化协议（1字节）
	readBuff = make([]byte, 0, SerializationCodeLength)
	for len(readBuff) < SerializationCodeLength {
		n := utils.AssertLength(reader.Read(tempBuff[:]))
		readBuff = append(readBuff, tempBuff[:n]...)
	}
	seria := GetSerializerByCode(BytesToInt(readBuff))
	utils.AssertIf(seria == nil, "Invalid serializer type:%X", readBuff)
	// 读取数据长度（4字节）
	readBuff = make([]byte, 0, DataLengthCodeLength)
	for len(readBuff) < DataLengthCodeLength {
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
	var mes any
	err := json.Unmarshal(dataBytes, &mes)
	if err != nil {
		log.Fatal("err: ", err)
	}

	data, err1 := seria.Deserialize(dataBytes, packageProtocol)

	if err1 != nil {
		log.Fatal("dataBytes deserialize failed：", err1)
	}

	// post message
	ctx.HandleRead(data)

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
