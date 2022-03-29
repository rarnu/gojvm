package gojvm

// 需要引用jni，先将jni相关的头文件链接到 /usr/local/include
// 所需文件为 jni.h, jni_md.h
// 同时需要将 libjvm.so/libjvm.dylib 链接到 /usr/local/lib

//#cgo CFLAGS: -I/usr/local/include
//#cgo LDFLAGS: -L/usr/local/lib -ljvm
/*
#include "gojvm_c.h"
*/
import "C"
import (
	"strings"
	"unsafe"
)

type JavaVM struct {
	jvm *C.JavaVM
}

type JavaEnv struct {
	jvm *C.JavaVM
	env *C.JNIEnv
}

type JavaClass struct {
	jvm       *C.JavaVM
	env       *C.JNIEnv
	clazz     C.jclass
	ClassName string
}

type JavaObject struct {
	jvm       *C.JavaVM
	env       *C.JNIEnv
	clazz     C.jclass
	obj       C.jobject
	ClassName string
}

//=============================================================
// jvm
//=============================================================

func NewJVM(classPath string, xms string, xmx string, xmn string, xss string) *JavaVM {
	cpath := C.CString(classPath)
	cxms := C.CString(xms)
	cxmx := C.CString(xmx)
	cxmn := C.CString(xmn)
	cxss := C.CString(xss)
	defer C.free(unsafe.Pointer(cpath))
	defer C.free(unsafe.Pointer(cxms))
	defer C.free(unsafe.Pointer(cxmx))
	defer C.free(unsafe.Pointer(cxmn))
	defer C.free(unsafe.Pointer(cxss))

	jvm := C.createJvm(cpath, cxms, cxmx, cxmn, cxss)
	if jvm == nil {
		return nil
	}

	return &JavaVM{jvm: jvm}
}

func (vm *JavaVM) Free() {
	_ = C.destroyJvm(vm.jvm)
}

func (vm *JavaVM) Attach() *JavaEnv {
	env := C.attachJvm(vm.jvm)
	if env == nil {
		return nil
	}
	return &JavaEnv{
		jvm: vm.jvm,
		env: env,
	}
}

//=============================================================
// env
//=============================================================

func (env *JavaEnv) Detach() {
	C.detachJvm(env.jvm)
}

func (env *JavaEnv) FindClass(className string) *JavaClass {
	cn := strings.ReplaceAll(className, ".", "/")
	cname := C.CString(cn)
	defer C.free(unsafe.Pointer(cname))
	clazz := C.findClass(env.env, cname)
	if clazz == 0 {
		return nil
	}
	return &JavaClass{
		jvm:       env.jvm,
		env:       env.env,
		clazz:     clazz,
		ClassName: className,
	}
}

func (env *JavaEnv) NewObject(className string) *JavaObject {
	jc := env.FindClass(className)
	if jc == nil {
		return nil
	}
	jo := C.newJavaObject(env.env, jc.clazz)
	if jo == C.jobject(C.NULL) {
		return nil
	}
	return &JavaObject{
		jvm:       env.jvm,
		env:       env.env,
		clazz:     jc.clazz,
		obj:       jo,
		ClassName: className,
	}
}

//=============================================================
// class
//=============================================================

func (c *JavaClass) InvokeVoid(methodName string, args ...any) error {
	types, err := ArgumentsCheck(args...)
	if err != nil {
		return err
	}
	sigStr, typArg, valArg := ParseArguments(types, "Ljava/lang/String;", args...)
	defer FreeArgs(len(args), types, typArg, valArg)
	cmn, csig := ParseNameSig(methodName, sigStr)
	defer FreeNameSig(cmn, csig)
	clen := C.int(len(args))
	C.callStaticVoidMethod(c.env, c.clazz, cmn, csig, clen, typArg, valArg)
	return nil
}

func (c *JavaClass) InvokeString(methodName string, args ...any) (string, error) {
	types, err := ArgumentsCheck(args...)
	if err != nil {
		return "", err
	}
	sigStr, typArg, valArg := ParseArguments(types, "Ljava/lang/String;", args...)
	defer FreeArgs(len(args), types, typArg, valArg)
	cmn, csig := ParseNameSig(methodName, sigStr)
	defer FreeNameSig(cmn, csig)
	clen := C.int(len(args))
	ret := C.callStaticStringMethod(c.env, c.clazz, cmn, csig, clen, typArg, valArg)
	if ret == nil {
		return "", nil
	}
	defer C.free(unsafe.Pointer(ret))
	return C.GoString(ret), nil
}

func (c *JavaClass) InvokeInt(methodName string, args ...any) (int, error) {
	types, err := ArgumentsCheck(args...)
	if err != nil {
		return 0, err
	}
	sigStr, typArg, valArg := ParseArguments(types, "I", args...)
	defer FreeArgs(len(args), types, typArg, valArg)
	cmn, csig := ParseNameSig(methodName, sigStr)
	defer FreeNameSig(cmn, csig)
	clen := C.int(len(args))
	ret := C.callStaticIntMethod(c.env, c.clazz, cmn, csig, clen, typArg, valArg)
	return int(ret), nil
}

func (c *JavaClass) InvokeLong(methodName string, args ...any) (int64, error) {
	types, err := ArgumentsCheck(args...)
	if err != nil {
		return 0, err
	}
	sigStr, typArg, valArg := ParseArguments(types, "I", args...)
	defer FreeArgs(len(args), types, typArg, valArg)
	cmn, csig := ParseNameSig(methodName, sigStr)
	defer FreeNameSig(cmn, csig)
	clen := C.int(len(args))
	ret := C.callStaticLongMethod(c.env, c.clazz, cmn, csig, clen, typArg, valArg)
	return int64(ret), nil
}

func (c *JavaClass) InvokeShort(methodName string, args ...any) (int16, error) {
	types, err := ArgumentsCheck(args...)
	if err != nil {
		return 0, err
	}
	sigStr, typArg, valArg := ParseArguments(types, "I", args...)
	defer FreeArgs(len(args), types, typArg, valArg)
	cmn, csig := ParseNameSig(methodName, sigStr)
	defer FreeNameSig(cmn, csig)
	clen := C.int(len(args))
	ret := C.callStaticShortMethod(c.env, c.clazz, cmn, csig, clen, typArg, valArg)
	return int16(ret), nil
}

func (c *JavaClass) InvokeByte(methodName string, args ...any) (uint8, error) {
	types, err := ArgumentsCheck(args...)
	if err != nil {
		return 0, err
	}
	sigStr, typArg, valArg := ParseArguments(types, "I", args...)
	defer FreeArgs(len(args), types, typArg, valArg)
	cmn, csig := ParseNameSig(methodName, sigStr)
	defer FreeNameSig(cmn, csig)
	clen := C.int(len(args))
	ret := C.callStaticByteMethod(c.env, c.clazz, cmn, csig, clen, typArg, valArg)
	return uint8(ret), nil
}

func (c *JavaClass) InvokeFloat(methodName string, args ...any) (float32, error) {
	types, err := ArgumentsCheck(args...)
	if err != nil {
		return 0, err
	}
	sigStr, typArg, valArg := ParseArguments(types, "I", args...)
	defer FreeArgs(len(args), types, typArg, valArg)
	cmn, csig := ParseNameSig(methodName, sigStr)
	defer FreeNameSig(cmn, csig)
	clen := C.int(len(args))
	ret := C.callStaticFloatMethod(c.env, c.clazz, cmn, csig, clen, typArg, valArg)
	return float32(ret), nil
}

func (c *JavaClass) InvokeDouble(methodName string, args ...any) (float64, error) {
	types, err := ArgumentsCheck(args...)
	if err != nil {
		return 0, err
	}
	sigStr, typArg, valArg := ParseArguments(types, "I", args...)
	defer FreeArgs(len(args), types, typArg, valArg)
	cmn, csig := ParseNameSig(methodName, sigStr)
	defer FreeNameSig(cmn, csig)
	clen := C.int(len(args))
	ret := C.callStaticDoubleMethod(c.env, c.clazz, cmn, csig, clen, typArg, valArg)
	return float64(ret), nil
}

func (c *JavaClass) InvokeBoolean(methodName string, args ...any) (bool, error) {
	types, err := ArgumentsCheck(args...)
	if err != nil {
		return false, err
	}
	sigStr, typArg, valArg := ParseArguments(types, "I", args...)
	defer FreeArgs(len(args), types, typArg, valArg)
	cmn, csig := ParseNameSig(methodName, sigStr)
	defer FreeNameSig(cmn, csig)
	clen := C.int(len(args))
	ret := C.callStaticBooleanMethod(c.env, c.clazz, cmn, csig, clen, typArg, valArg)
	return int(ret) != 0, nil
}

func (c *JavaClass) GetObject(fieldName string, className string) *JavaObject {
	sig := strings.ReplaceAll(className, ".", "/")
	cmn, csig := ParseNameSig(fieldName, "L"+sig+";")
	defer FreeNameSig(cmn, csig)
	ret := C.getStaticObject(c.env, c.clazz, cmn, csig)
	if ret == C.jobject(C.NULL) {
		return nil
	}
	return &JavaObject{
		jvm:       c.jvm,
		env:       c.env,
		clazz:     C.findClass(c.env, C.CString(sig)),
		obj:       ret,
		ClassName: className,
	}
}

func (c *JavaClass) SetObject(fieldName string, className string, obj *JavaObject) {
	sig := strings.ReplaceAll(className, ".", "/")
	cmn, csig := ParseNameSig(fieldName, "L"+sig+";")
	defer FreeNameSig(cmn, csig)
	C.setStaticObject(c.env, c.clazz, cmn, csig, obj.obj)
}

func (c *JavaClass) GetString(fieldName string) string {
	cmn := C.CString(fieldName)
	defer C.free(unsafe.Pointer(cmn))
	ret := C.getStaticString(c.env, c.clazz, cmn)
	if ret == nil {
		return ""
	}
	defer C.free(unsafe.Pointer(ret))
	return C.GoString(ret)
}

func (c *JavaClass) SetString(fieldName string, value string) {
	cmn := C.CString(fieldName)
	defer C.free(unsafe.Pointer(cmn))
	cval := C.CString(value)
	defer C.free(unsafe.Pointer(cval))
	C.setStaticString(c.env, c.clazz, cmn, cval)
}

func (c *JavaClass) GetInt(fieldName string) int {
	cmn := C.CString(fieldName)
	defer C.free(unsafe.Pointer(cmn))
	ret := C.getStaticInt(c.env, c.clazz, cmn)
	return int(ret)
}

func (c *JavaClass) SetInt(fieldName string, value int) {
	cmn := C.CString(fieldName)
	defer C.free(unsafe.Pointer(cmn))
	C.setStaticInt(c.env, c.clazz, cmn, C.int(value))
}

func (c *JavaClass) Free() {
	C.freeJavaClassRef(c.env, c.clazz)
}

//=============================================================
// class
//=============================================================

func (o *JavaObject) Free() {
	C.freeJavaClassRef(o.env, o.clazz)
	C.freeJavaObject(o.env, o.obj)
}

func (o *JavaObject) GetString(fieldName string) string {
	cmn := C.CString(fieldName)
	defer C.free(unsafe.Pointer(cmn))
	ret := C.getObjString(o.env, o.clazz, o.obj, cmn)
	if ret == nil {
		return ""
	}
	defer C.free(unsafe.Pointer(ret))
	return C.GoString(ret)
}

func (o *JavaObject) SetString(fieleName string, value string) {
	cmn := C.CString(fieleName)
	defer C.free(unsafe.Pointer(cmn))
	cval := C.CString(value)
	defer C.free(unsafe.Pointer(cval))
	C.setObjString(o.env, o.clazz, o.obj, cmn, cval)
}
