#ifndef __GOJVM_H__
#define __GOJVM_H__

#include <stdio.h>
#include <stdlib.h>
#include <jni.h>
#include <string.h>
#include <stdbool.h>

JavaVM* createJvm(char* classPath, char* xms, char* xmx, char* xmn, char* xss);
int destroyJvm(JavaVM* jvm);
int destroyJvm(JavaVM* jvm);
JNIEnv* attachJvm(JavaVM* jvm);
void detachJvm(JavaVM* jvm);
jclass findClass(JNIEnv* env, char* className);

void callStaticVoidMethod(JNIEnv* env, jclass clazz, char* methodName, char* sig, int len, char** types, void** args);
char* callStaticStringMethod(JNIEnv* env, jclass clazz, char* methodName, char* sig, int len, char** types, void** args);
int callStaticIntMethod(JNIEnv* env, jclass clazz, char* methodName, char* sig, int len, char** types, void** args);
long callStaticLongMethod(JNIEnv* env, jclass clazz, char* methodName, char* sig, int len, char** types, void** args);
short callStaticShortMethod(JNIEnv* env, jclass clazz, char* methodName, char* sig, int len, char** types, void** args);
unsigned char callStaticByteMethod(JNIEnv* env, jclass clazz, char* methodName, char* sig, int len, char** types, void** args);
float callStaticFloatMethod(JNIEnv* env, jclass clazz, char* methodName, char* sig, int len, char** types, void** args);
double callStaticDoubleMethod(JNIEnv* env, jclass clazz, char* methodName, char* sig, int len, char** types, void** args);
int callStaticBooleanMethod(JNIEnv* env, jclass clazz, char* methodName, char* sig, int len, char** types, void** args);

jobject getStaticObject(JNIEnv* env, jclass clazz, char* fieldName, char* sig);
void setStaticObject(JNIEnv* env, jclass clazz, char* fieldName, char* sig, jobject obj);
char* getStaticString(JNIEnv* env, jclass clazz, char* fieldName);
void setStaticString(JNIEnv* env, jclass clazz, char* fieldName, char* value);
int getStaticInt(JNIEnv* env, jclass clazz, char* fieldName);
void setStaticInt(JNIEnv* env, jclass clazz, char* fieldName, int value);


char* getObjString(JNIEnv* env, jclass clazz, jobject obj, char* fieldName);
void setObjString(JNIEnv* env, jclass clazz, jobject obj, char* fieldName, char* value);


jobject newJavaObject(JNIEnv* env, jclass clazz);
void freeJavaClassRef(JNIEnv* env, jclass clz);
void freeJavaObject(JNIEnv* env, jobject obj);

#endif // __GOJVM_H__