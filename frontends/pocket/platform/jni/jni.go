package jni

/*
#include <jni.h>
#include <dlfcn.h>
#include <stdlib.h>

JavaVM *vm = NULL;
jobject ctx = NULL;

void initialize_pointers() {
	vm = *(JavaVM**)dlsym(RTLD_DEFAULT, "current_vm");
	ctx = *(jobject*)dlsym(RTLD_DEFAULT, "current_ctx");
}

int test() {
	JNIEnv *env;
	if ((*vm)->GetEnv(vm, (void**)&env, JNI_VERSION_1_6) != JNI_OK) {
		return -1;
	}

	jclass clazz = (*env)->GetObjectClass(env, ctx);
	if (clazz == NULL) {
		return -2;
	}

    jmethodID getIntentMethodID = (*env)->GetMethodID(env, clazz, "getIntent", "()Landroid/content/Intent;");
	if (getIntentMethodID == NULL) {
		return -3;
	}

	jobject intent = (*env)->CallObjectMethod(env, clazz, getIntentMethodID);
	if (intent == NULL) {
		return -4;
	}


	return 0;
}
*/
import "C"

import (
	"log"
	"unsafe"
)

func Test() unsafe.Pointer {
	C.initialize_pointers()
	if C.vm == nil || C.ctx == nil {
		log.Fatalln("Could not initialize JVM references")
	}

	log.Println("aaaaaaaaaaaaaaaaaaaaaaaaaa")

	log.Println("##########", C.test())
	return nil
}
