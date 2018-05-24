package jni

/*
#include <jni.h>
#include <dlfcn.h>
#include <stdlib.h>

JavaVM *vm = NULL;
jobject ctx = NULL;

static void initialize_pointers() {
	if (vm == NULL)
		vm = *(JavaVM**)dlsym(RTLD_DEFAULT, "current_vm");
	if (ctx == NULL)
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
	if (clazz == NULL)
		return -1;

    jmethodID getIntentMethodID = (*env)->GetMethodID(env, clazz, "getIntent", "()Landroid/content/Intent;");
	if (getIntentMethodID == NULL)
		return -2;

	jobject intent = (*env)->CallObjectMethod(env, ctx, getIntentMethodID);
	if (intent == NULL)
		return -3;

	clazz = (*env)->GetObjectClass(env, intent);
	if (clazz == NULL)
		return -4;

	jmethodID getDataMethodID = (*env)->GetMethodID(env, clazz, "getData", "()Landroid/net/Uri;");
	if (getDataMethodID == NULL)
		return -5;

	jobject uri = (*env)->CallObjectMethod(env, intent, getDataMethodID);
	if (uri == NULL)
		return -6;

	clazz = (*env)->GetObjectClass(env, uri);
	if (clazz == NULL)
		return -7;

	jmethodID toStringMethodID = (*env)->GetMethodID(env, clazz, "toString", "()Ljava/lang/String;");
	if (toStringMethodID == NULL)
		return -8;

	jstring jstr = (*env)->CallObjectMethod(env, uri, toStringMethodID);
	if (intent == NULL)
		return -9;

	*url = (*env)->GetStringUTFChars(env, jstr, NULL);
	if (*url == NULL)
		return -10;

	return 0;
}

static int openURL(uintptr_t e, const char *url) {
	JNIEnv *env = (JNIEnv*)e;

	jstring str = (*env)->NewStringUTF(env, url);
	if (str == NULL)
		return -1;

	jclass clazz = (*env)->FindClass(env, "android/net/Uri");
	if (clazz == NULL)
		return -2;

	jmethodID parseMethodID = (*env)->GetStaticMethodID(env, clazz, "parse", "(Ljava/lang/String;)Landroid/net/Uri;");
	if (parseMethodID == NULL)
		return -3;

	jobject uri = (*env)->CallStaticObjectMethod(env, clazz, parseMethodID, str);
	if (uri == NULL)
		return -4;

	clazz = (*env)->FindClass(env, "android/content/Intent");
	if (clazz == NULL)
		return -5;

	jfieldID actionView = (*env)->GetStaticFieldID(env, clazz, "ACTION_VIEW", "Ljava/lang/String;");
	if (actionView == NULL)
		return -6;

	jobject actionViewObject = (*env)->GetStaticObjectField(env, clazz, actionView);
	if (actionViewObject == NULL)
		return -7;

    jmethodID intentConstructor = (*env)->GetMethodID(env, clazz, "<init>", "(Ljava/lang/String;Landroid/net/Uri;)V");
	if (intentConstructor == NULL)
		return -8;

	jobject intent = (*env)->NewObject(env, clazz, intentConstructor, actionViewObject, uri);
	if (intent == NULL)
		return -9;

	clazz = (*env)->GetObjectClass(env, ctx);
	if (clazz == NULL)
		return -10;

    jmethodID startActivityMethodID = (*env)->GetMethodID(env, clazz, "startActivity", "(Landroid/content/Intent;)V");
	if (startActivityMethodID == NULL)
		return -11;

	(*env)->CallVoidMethod(env, ctx, startActivityMethodID, intent);
	return 0;
}
*/
import "C"

import (
	"errors"
	"fmt"
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
		return "", fmt.Errorf("could not read URL: %v", ret)
	}
	return C.GoString(str), nil
}

func OpenURL(url string) error {
	C.initialize_pointers()
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	env := C.uintptr_t(0)
	attached := C.int(0)
	if errStr := C.lockJNI(&env, &attached); errStr != nil {
		return errors.New(C.GoString(errStr))
	}
	if attached != 0 {
		defer C.unlockJNI()
	}

	if ret := C.openURL(env, C.CString(url)); ret != 0 {
		return fmt.Errorf("could not open URL: %v", ret)
	}
	return nil
}
