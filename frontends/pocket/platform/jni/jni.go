package jni

/*
#include <jni.h>
#include <dlfcn.h>
#include <stdlib.h>

JavaVM *vm = NULL;
jobject ctx = NULL;

static void initialize_pointers() {
	vm = *(JavaVM**)dlsym(RTLD_DEFAULT, "current_vm");
	ctx = *(jobject*)dlsym(RTLD_DEFAULT, "current_ctx");
}

static char* lockJNI(uintptr_t* envp, int* attachedp) {
	JNIEnv* env;

	if (vm == NULL) {
		return "no current JVM";
	}

	*attachedp = 0;
	switch ((*vm)->GetEnv(vm, (void**)&env, JNI_VERSION_1_6)) {
	case JNI_OK:
		break;
	case JNI_EDETACHED:
		if ((*vm)->AttachCurrentThread(vm, &env, 0) != 0) {
			return "cannot attach to JVM";
		}
		*attachedp = 1;
		break;
	case JNI_EVERSION:
		return "bad JNI version";
	default:
		return "unknown JNI error from GetEnv";
	}

	*envp = (uintptr_t)env;
	return NULL;
}

static void unlockJNI() {
	(*vm)->DetachCurrentThread(vm);
}

static int getURL(uintptr_t e, const char **url) {
	JNIEnv *env = (JNIEnv*)e;

	jclass clazz = (*env)->GetObjectClass(env, ctx);
	if (clazz == NULL) {
		return -1;
	}

    jmethodID getIntentMethodID = (*env)->GetMethodID(env, clazz, "getIntent", "()Landroid/content/Intent;");
	if (getIntentMethodID == NULL) {
		return -2;
	}

	jobject intent = (*env)->CallObjectMethod(env, ctx, getIntentMethodID, NULL);
	if (intent == NULL) {
		return -3;
	}

	clazz = (*env)->GetObjectClass(env, intent);
	if (clazz == NULL) {
		return -4;
	}

	jmethodID getDataMethodID = (*env)->GetMethodID(env, clazz, "getData", "()Landroid/net/Uri;");
	if (getDataMethodID == NULL) {
		return -5;
	}

	jobject uri = (*env)->CallObjectMethod(env, intent, getDataMethodID, NULL);
	if (uri == NULL) {
		return -6;
	}

	clazz = (*env)->GetObjectClass(env, uri);
	if (clazz == NULL) {
		return -7;
	}

	jmethodID toStringMethodID = (*env)->GetMethodID(env, clazz, "toString", "()Ljava/lang/String;");
	if (toStringMethodID == NULL) {
		return -8;
	}

	jstring jstr = (*env)->CallObjectMethod(env, uri, toStringMethodID, NULL);
	if (intent == NULL) {
		return -9;
	}

	*url = (*env)->GetStringUTFChars(env, jstr, NULL);
	if (*url == NULL) {
		return -10;
	}

	return 0;
}
*/
import "C"

import (
	"errors"
	"runtime"
)

func GetURL() (string, error) {
	C.initialize_pointers()
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	env := C.uintptr_t(0)
	attached := C.int(0)
	if errStr := C.lockJNI(&env, &attached); errStr != nil {
		return "", errors.New(C.GoString(errStr))
	}
	if attached != 0 {
		defer C.unlockJNI()
	}

	var str *C.char
	if ret := C.getURL(env, &str); ret != 0 {
		return "", errors.New("could not read URL")
	}
	return C.GoString(str), nil
}
