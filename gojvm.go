package gojvm

// 需要引用jni，先将jni相关的头文件链接到 /usr/local/include
// 所需文件为 jni.h, jni_md.h
// 同时需要将 libjvm.so/libjvm.dylib 链接到 /usr/local/lib

//#cgo CFLAGS: -I/usr/local/include
//#cgo LDFLAGS: -L/usr/local/lib -ljvm
/*
#include <stdio.h>
#include <stdlib.h>
#include <jni.h>
#include <string.h>
#include <stdbool.h>

JavaVM* createJvm(char* classPath, char* xms, char* xmx, char* xmn, char* xss) {
	JavaVM* jvm;
	JNIEnv* env;
	JavaVMInitArgs vm_args;
	JavaVMOption options[5];

	options[0].optionString = (char*)malloc(strlen("-Djava.class.path=") + strlen(classPath) + 1);
	sprintf(options[0].optionString, "-Djava.class.path=%s", classPath);
	options[1].optionString = (char*)malloc(strlen("-Xms") + strlen(xms) + 1);
	sprintf(options[1].optionString, "-Xms%s", xms);
	options[2].optionString = (char*)malloc(strlen("-Xmx") + strlen(xmx) + 1);
	sprintf(options[2].optionString, "-Xmx%s", xmx);
	options[3].optionString = (char*)malloc(strlen("-Xmn") + strlen(xmn) + 1);
	sprintf(options[3].optionString, "-Xmn%s", xmn);
	options[4].optionString = (char*)malloc(strlen("-Xss") + strlen(xss) + 1);
	sprintf(options[4].optionString, "-Xss%s", xss);

	vm_args.version = JNI_VERSION_1_8;
	vm_args.nOptions = 5;
	vm_args.options = options;
	vm_args.ignoreUnrecognized = JNI_FALSE;

	jint res = JNI_CreateJavaVM(&jvm, (void**)&env, &vm_args);
	if (res < 0) {
		printf("create jvm failed\n");
		return NULL;
	}
	(*jvm)->DetachCurrentThread(jvm);
	return jvm;
}

int destroyJvm(JavaVM* jvm) {
	jint res = (*jvm)->DestroyJavaVM(jvm);
	if (res < 0) {
		printf("destroy jvm failed\n");
		return 1;
	}
	return 0;
}

JNIEnv* attachJvm(JavaVM* jvm) {
	JNIEnv* env;
	jint res = (*jvm)->AttachCurrentThread(jvm, (void**)&env, NULL);
	if (res < 0) {
		printf("attach jvm failed\n");
		return NULL;
	}
	return env;
}

void detachJvm(JavaVM* jvm) {
	(*jvm)->DetachCurrentThread(jvm);
}

jclass findClass(JNIEnv* env, char* className) {
	jclass cls = (*env)->FindClass(env, className);
	if (cls == NULL) {
		printf("find class failed\n");
		return NULL;
	}
	return cls;
}

jvalue* makeParams(JNIEnv* env, int len, char** types, void** args) {
	jvalue *v = malloc(sizeof(jvalue) * len);
	for (int i = 0; i < len; i++) {
		if (strcmp(types[i], "Ljava/lang/String;") == 0) {
			v[i].l = (*env)->NewStringUTF(env, (char*)args[i]);
		} else if (strcmp(types[i], "I") == 0) {
			v[i].i = *((int*)args[i]);
		}
	}
	return v;
}

void freeParams(JNIEnv* env, int len, char** types, jvalue* v) {
	for (int i = 0; i < len; i++) {
		if (strcmp(types[i], "Ljava/lang/String;") == 0) {
			(*env)->DeleteLocalRef(env, v[i].l);
		}
	}
	free(v);
}

void callStaticVoidMethod(JNIEnv* env, jclass clazz, char* methodName, char* sig, int len, char** types, void** args) {
	jmethodID m = (*env)->GetStaticMethodID(env, clazz, methodName, sig);
	jvalue *v = makeParams(env, len, types, args);
	(*env)->CallStaticVoidMethodA(env, clazz, m, v);
	freeParams(env, len, types, v);
}

jobject callStaticObjectMethod(JNIEnv* env, jclass clazz, char* methodName, char* sig, int len, char** types, void** args) {
	jmethodID m = (*env)->GetStaticMethodID(env, clazz, methodName, sig);
	jvalue *v = makeParams(env, len, types, args);
	jobject jobj = (*env)->CallStaticObjectMethodA(env, clazz, m, v);
	freeParams(env, len, types, v);
	return jobj;
}

char* callStaticStringMethod(JNIEnv* env, jclass clazz, char* methodName, char* sig, int len, char** types, void** args) {
	jstring jstr = callStaticObjectMethod(env, clazz, methodName, sig, len, types, args);
	const char* str = (*env)->GetStringUTFChars(env, jstr, NULL);
	(*env)->DeleteLocalRef(env, jstr);
	return (char*)str;
}


*/
import "C"
import (
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
	cname := C.CString(className)
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

//=============================================================
// class
//=============================================================

func (c *JavaClass) CallStaticVoidMethod(methodName string, args ...any) error {
	types, err := ArgumentsCheck(args...)
	if err != nil {
		return err
	}
	sigStr, typArg, valArg, _ := ParseArguments(types, "Ljava/lang/String;", args...)
	defer FreeArgs(len(args), types, typArg, valArg)
	cmn := C.CString(methodName)
	defer C.free(unsafe.Pointer(cmn))
	csig := C.CString(sigStr)
	defer C.free(unsafe.Pointer(csig))
	clen := C.int(len(args))
	C.callStaticVoidMethod(c.env, c.clazz, cmn, csig, clen, typArg, valArg)
	return nil
}

func (c *JavaClass) CallStaticStringMethod(methodName string, args ...any) (string, error) {
	types, err := ArgumentsCheck(args...)
	if err != nil {
		return "", err
	}
	sigStr, typArg, valArg, _ := ParseArguments(types, "Ljava/lang/String;", args...)
	defer FreeArgs(len(args), types, typArg, valArg)
	cmn := C.CString(methodName)
	defer C.free(unsafe.Pointer(cmn))
	csig := C.CString(sigStr)
	defer C.free(unsafe.Pointer(csig))
	clen := C.int(len(args))
	ret := C.callStaticStringMethod(c.env, c.clazz, cmn, csig, clen, typArg, valArg)
	defer C.free(unsafe.Pointer(ret))
	return C.GoString(ret), nil
}

func (c *JavaClass) CallStaticSliceMethod(methodName string, args ...any) []any {
	return []any{}
}

func (c *JavaClass) CallStaticMapMethod(methodName string, args ...any) map[string]any {
	return map[string]any{}
}
