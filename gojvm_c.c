
#include "gojvm_c.h"

bool hasPrefix(const char *str, const char *sub) {
    return strncmp(str, sub, strlen(sub)) == 0;
}

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


char* getSubType(char* type) {
	size_t len = strlen(type);
    char* s = (char*)malloc(len);
    int inType = 0;
    int iidx = 0;
    for (int i = 0; i < len; i++) {
        if (type[i] == '<') {
            inType = 1;
            continue;
        }
        if (type[i] == '>') break;
        if (inType) s[iidx++] = type[i];
    }
    return s;
}

char* getRealSig(char* sig) {
	size_t len = strlen(sig);
    char* s = (char*) malloc(len);
    int inType = 0;
    int iidx = 0;
    for (int i = 0; i < len; i++) {
        if (sig[i] == '<') {
            inType = 1;
            continue;
        }
        if (sig[i] == '>') {
            inType = 0;
            continue;
        }
        if (!inType) {
            s[iidx++] = sig[i];
        }
    }
    return s;
}

jvalue* makeParams(JNIEnv* env, int len, char** types, void** args) {
	jvalue *v = malloc(sizeof(jvalue) * len);
	for (int i = 0; i < len; i++) {
		if (strcmp(types[i], "Ljava/lang/String;") == 0) {
			v[i].l = (*env)->NewStringUTF(env, (char*)args[i]);
		} else if (strcmp(types[i], "I") == 0) {
			v[i].i = *((int*)args[i]);
		} else if (strcmp(types[i], "J") == 0) {
			v[i].j = *((long*)args[i]);
		} else if (strcmp(types[i], "F") == 0) {
            v[i].f = *((float*)args[i]);
        } else if (strcmp(types[i], "D") == 0) {
            v[i].d = *((double*)args[i]);
        } else if (strcmp(types[i], "B") == 0) {
            v[i].b = *((unsigned char*)args[i]);
        } else if (strcmp(types[i], "S") == 0) {
            v[i].s = *((short*)args[i]);
        } else if (strcmp(types[i], "Z") == 0) {
            int bi = *((int*)args[i]);
            v[i].z = bi == 0 ? JNI_FALSE : JNI_TRUE;
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

int callStaticIntMethod(JNIEnv* env, jclass clazz, char* methodName, char* sig, int len, char** types, void** args) {
    jmethodID m = (*env)->GetStaticMethodID(env, clazz, methodName, sig);
    jvalue *v = makeParams(env, len, types, args);
    jint jret = (*env)->CallStaticIntMethodA(env, clazz, m, v);
    freeParams(env, len, types, v);
    return (int)jret;
}

long callStaticLongMethod(JNIEnv* env, jclass clazz, char* methodName, char* sig, int len, char** types, void** args) {
    jmethodID m = (*env)->GetStaticMethodID(env, clazz, methodName, sig);
    jvalue *v = makeParams(env, len, types, args);
    jlong jret = (*env)->CallStaticLongMethodA(env, clazz, m, v);
    freeParams(env, len, types, v);
    return (long)jret;
}

short callStaticShortMethod(JNIEnv* env, jclass clazz, char* methodName, char* sig, int len, char** types, void** args) {
    jmethodID m = (*env)->GetStaticMethodID(env, clazz, methodName, sig);
    jvalue *v = makeParams(env, len, types, args);
    jshort jret = (*env)->CallStaticShortMethodA(env, clazz, m, v);
    freeParams(env, len, types, v);
    return (short)jret;
}

unsigned char callStaticByteMethod(JNIEnv* env, jclass clazz, char* methodName, char* sig, int len, char** types, void** args) {
    jmethodID m = (*env)->GetStaticMethodID(env, clazz, methodName, sig);
    jvalue *v = makeParams(env, len, types, args);
    jbyte jret = (*env)->CallStaticByteMethodA(env, clazz, m, v);
    freeParams(env, len, types, v);
    return (unsigned char)jret;
}

float callStaticFloatMethod(JNIEnv* env, jclass clazz, char* methodName, char* sig, int len, char** types, void** args) {
    jmethodID m = (*env)->GetStaticMethodID(env, clazz, methodName, sig);
    jvalue *v = makeParams(env, len, types, args);
    jfloat jret = (*env)->CallStaticFloatMethodA(env, clazz, m, v);
    freeParams(env, len, types, v);
    return (float)jret;
}

double callStaticDoubleMethod(JNIEnv* env, jclass clazz, char* methodName, char* sig, int len, char** types, void** args) {
    jmethodID m = (*env)->GetStaticMethodID(env, clazz, methodName, sig);
    jvalue *v = makeParams(env, len, types, args);
    jfloat jret = (*env)->CallStaticFloatMethodA(env, clazz, m, v);
    freeParams(env, len, types, v);
    return (float)jret;
}

int callStaticBooleanMethod(JNIEnv* env, jclass clazz, char* methodName, char* sig, int len, char** types, void** args) {
    jmethodID m = (*env)->GetStaticMethodID(env, clazz, methodName, sig);
    jvalue *v = makeParams(env, len, types, args);
    jboolean jret = (*env)->CallStaticBooleanMethodA(env, clazz, m, v);
    freeParams(env, len, types, v);
    return (int)jret;
}

jobject getStaticObject(JNIEnv* env, jclass clazz, char* fieldName, char* sig) {
    jfieldID f = (*env)->GetStaticFieldID(env, clazz, fieldName, sig);
    jobject jret = (*env)->GetStaticObjectField(env, clazz, f);
    return jret;
}

void setStaticObject(JNIEnv* env, jclass clazz, char* fieldName, char* sig, jobject obj) {
    jfieldID f = (*env)->GetStaticFieldID(env, clazz, fieldName, sig);
    (*env)->SetStaticObjectField(env, clazz, f, obj);
}

char* getStaticString(JNIEnv* env, jclass clazz, char* fieldName) {
    jfieldID f = (*env)->GetStaticFieldID(env, clazz, fieldName, "Ljava/lang/String;");
    jstring jret = (*env)->GetStaticObjectField(env, clazz, f);
    const char* cret = (*env)->GetStringUTFChars(env, jret, NULL);
    (*env)->DeleteLocalRef(env, jret);
    return (char*)cret;
}

void setStaticString(JNIEnv* env, jclass clazz, char* fieldName, char* value) {
    jfieldID f = (*env)->GetStaticFieldID(env, clazz, fieldName, "Ljava/lang/String;");
    jstring jstr = (*env)->NewStringUTF(env, value);
    (*env)->SetStaticObjectField(env, clazz, f, jstr);
    (*env)->DeleteLocalRef(env, jstr);
}

int getStaticInt(JNIEnv* env, jclass clazz, char* fieldName) {
    jfieldID f = (*env)->GetStaticFieldID(env, clazz, fieldName, "I");
    jint jret = (*env)->GetStaticIntField(env, clazz, f);
    return (int)jret;
}

void setStaticInt(JNIEnv* env, jclass clazz, char* fieldName, int value) {
    jfieldID f = (*env)->GetStaticFieldID(env, clazz, fieldName, "I");
    (*env)->SetStaticIntField(env, clazz, f, (jint)value);
}




char* getObjString(JNIEnv* env, jclass clazz, jobject obj, char* fieldName) {
    jfieldID f = (*env)->GetFieldID(env, clazz, fieldName, "Ljava/lang/String;");
    jstring jret = (*env)->GetObjectField(env, obj, f);
    const char* cret = (*env)->GetStringUTFChars(env, jret, NULL);
    (*env)->DeleteLocalRef(env, jret);
    return (char*)cret;
}

void setObjString(JNIEnv* env, jclass clazz, jobject obj, char* fieldName, char* value) {
    jfieldID f = (*env)->GetFieldID(env, clazz, fieldName, "Ljava/lang/String;");
    jstring jstr = (*env)->NewStringUTF(env, value);
    (*env)->SetObjectField(env, obj, f, jstr);
    (*env)->DeleteLocalRef(env, jstr);
}

jobject newJavaObject(JNIEnv* env, jclass clazz) {
    jmethodID m = (*env)->GetMethodID(env, clazz, "<init>", "()V");
    jobject jret = (*env)->NewObject(env, clazz, m);
    return jret;
}

void freeJavaObject(JNIEnv* env, jobject obj) {
    (*env)->DeleteLocalRef(env, obj);
}

void freeJavaClassRef(JNIEnv* env, jclass clz) {
    (*env)->DeleteLocalRef(env, clz);
}

