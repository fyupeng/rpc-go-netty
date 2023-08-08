package rpc_go_netty

//
//import (
//	"encoding/json"
//	"fmt"
//	"rpc-go-netty/codec"
//	"io"
//	"log"
//	"net"
//	"reflect"
//	"sync"
//)
//
//const MagicNumber = 0x3bef5c
//
//type Option struct {
//	MagicNumber int        // MagicNumber marks
//	CodecType   codec.Type // Codec way
//}
//
//var DefaultOption = &Option{
//	MagicNumber: MagicNumber,
//	CodecType:   codec.GobType,
//}
//
//type Server struct{}
//
//func NewServer() *Server {
//	return &Server{}
//}
//
//var DefaultServer = NewServer()
//
//func (server *Server) Accept(lis net.Listener) {
//	for {
//		conn, err := lis.Accept()
//		if err != nil {
//			log.Println("rpc server: accept error:", err)
//			return
//		}
//		go server.ServerConn(conn)
//	}
//}
//
//func Accept(lis net.Listener) {
//	DefaultServer.Accept(lis)
//}
//
//func (server *Server) ServerConn(conn io.ReadWriteCloser) {
//	defer func() { _ = conn.Close() }()
//	var opt Option
//	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
//		log.Println("rpc server:options error:", err)
//		return
//	}
//
//	if opt.MagicNumber != MagicNumber {
//		log.Println("rpc server: invalid magic number %x", opt.MagicNumber)
//		return
//	}
//	f := codec.NewCodecFuncMap[opt.CodecType]
//	if f == nil {
//		log.Println("rpc server: invalid codec type %s", opt.CodecType)
//		return
//	}
//	server.serveCodec(f(conn))
//}
//
//var invalidRequest = struct{}{}
//
//func (server *Server) serveCodec(cc codec.Codec) {
//	sending := new(sync.Mutex) // make sure to send a complete response
//	wg := new(sync.WaitGroup)  // wait until all request are handled
//	for {
//		req, err := server.readRequest(cc)
//		if err != nil {
//			if req == nil {
//				break // it's not possible to recover, so close the connection
//			}
//			req.h.Error = err.Error()
//			server.sendResponse(cc, req.h, invalidRequest, sending)
//		}
//		wg.Add(1)
//		go server.handleRequest(cc, req, sending, wg)
//	}
//	wg.Wait()
//	_ = cc.Close()
//}
//
//type request struct {
//	h            *codec.Header // header of request
//	argv, replyv reflect.Value //argv and reply of request
//}
//
//func (server *Server) readRequestHeader(cc codec.Codec) (*codec.Header, error) {
//	var h codec.Header
//	if err := cc.ReadHeader(&h); err != nil {
//		if err != io.EOF && err != io.ErrUnexpectedEOF {
//			log.Println("rpc server: read header error:", err)
//		}
//		return nil, err
//	}
//	return &h, nil
//}
//
//func (server *Server) readRequest(cc codec.Codec) (*request, error) {
//	h, err := server.readRequestHeader(cc)
//	if err != nil {
//		return nil, err
//	}
//	req := &request{h: h}
//	// just suppose it's string
//	req.argv = reflect.New(reflect.TypeOf(""))
//	if err = cc.ReadBody(req.argv.Interface()); err != nil {
//		log.Println("rpc server: read argv err:", err)
//	}
//	return req, nil
//}
//
//func (Server *Server) sendResponse(cc codec.Codec, h *codec.Header, body interface{}, sending *sync.Mutex) {
//	// 加互斥锁/等待
//	sending.Lock()
//	// 延迟释放锁
//	defer sending.Unlock()
//	//
//	if err := cc.Write(h, body); err != nil {
//		log.Println("rpc server: write response error:", err)
//	}
//}
//
//func (server *Server) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {
//	// should call registered rpc methods to get the right replyv
//	// just print argv and send a hello message
//	defer wg.Done()
//	log.Println(req.h, req.argv.Elem())
//	req.replyv = reflect.ValueOf(fmt.Sprintf("geerpc resp %d", req.h.Seq))
//	server.sendResponse(cc, req.h, req.replyv.Interface(), sending)
//}