package ibmmq

/*
  Copyright (c) IBM Corporation 2016

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific

   Contributors:
     Mark Taylor - Initial Contribution
*/

/*
#cgo CFLAGS: -I/opt/mqm/inc
#cgo LDFLAGS: -L/opt/mqm/lib64 -lmqm -m64

#include <stdlib.h>
#include <string.h>
#include <cmqc.h>

*/
import "C"

import (
	"bytes"
	"unsafe"
)

/*
 * This file contains operations on the MQ Object Descriptor (MQOD)
 *
 */

type MQOD struct {
	StrucId         string
	Version         C.MQLONG
	ObjectType      C.MQLONG
	ObjectName      string
	ObjectQMgrName  string
	DynamicQName    string
	AlternateUserId string

	RecsPresent       C.MQLONG
	KnownDestCount    C.MQLONG
	UnknownDestCount  C.MQLONG
	InvalidDestCount  C.MQLONG
	ObjectRecOffset   C.MQLONG
	ResponseRecOffset C.MQLONG

	ObjectRecPtr   C.MQPTR
	ResponseRecPtr C.MQPTR

	AlternateSecurityId []byte
	ResolvedQName       string
	ResolvedQMgrName    string

	ObjectString    string
	SelectionString string
	ResObjectString string
	ResolvedType    C.MQLONG
}

func NewMQOD() *MQOD {

	od := new(MQOD)
	od.StrucId = "OD  "
	od.Version = 1
	od.ObjectType = C.MQOT_Q
	od.ObjectName = ""
	od.ObjectQMgrName = ""
	od.DynamicQName = "AMQ.*"
	od.AlternateUserId = ""

	od.RecsPresent = 0
	od.KnownDestCount = 0
	od.UnknownDestCount = 0
	od.InvalidDestCount = 0
	od.ObjectRecOffset = 0
	od.ResponseRecOffset = 0

	od.ObjectRecPtr = nil
	od.ResponseRecPtr = nil

	od.AlternateSecurityId = bytes.Repeat([]byte{0}, C.MQ_SECURITY_ID_LENGTH)
	od.ResolvedQName = ""
	od.ResolvedQMgrName = ""

	od.ObjectString = ""
	od.SelectionString = ""
	od.ResObjectString = ""
	od.ResolvedType = C.MQOT_NONE
	return od
}

/*
 * It is expected that copyXXtoC and copyXXfromC will be called as
 * matching pairs. That means that we can handle the MQCHARV type
 * by allocating storage in the toC function and freeing it in fromC.
 * If the input string for an MQCHARV type is empty, then we allocate
 * a fixed length buffer for any potential output.
 *
 * In the fromC function, that buffer is freed. Conveniently, we can
 * free it always, because if we didn't explicitly call malloc(), it was
 * allocated by C.CString and still needs to be freed.
 */
func copyODtoC(mqod *C.MQOD, good *MQOD) {
	var i int
	const vsbufsize = 10240
	setMQIString((*C.char)(&mqod.StrucId[0]), good.StrucId, 4)
	mqod.Version = good.Version
	mqod.ObjectType = good.ObjectType
	setMQIString((*C.char)(&mqod.ObjectName[0]), good.ObjectName, C.MQ_OBJECT_NAME_LENGTH)
	setMQIString((*C.char)(&mqod.ObjectQMgrName[0]), good.ObjectQMgrName, C.MQ_OBJECT_NAME_LENGTH)
	setMQIString((*C.char)(&mqod.DynamicQName[0]), good.DynamicQName, C.MQ_OBJECT_NAME_LENGTH)
	setMQIString((*C.char)(&mqod.AlternateUserId[0]), good.AlternateUserId, C.MQ_USER_ID_LENGTH)

	mqod.RecsPresent = good.RecsPresent
	mqod.KnownDestCount = good.KnownDestCount
	mqod.UnknownDestCount = good.UnknownDestCount
	mqod.InvalidDestCount = good.InvalidDestCount
	mqod.ObjectRecOffset = good.ObjectRecOffset
	mqod.ResponseRecOffset = good.ResponseRecOffset

	mqod.ObjectRecPtr = good.ObjectRecPtr
	mqod.ResponseRecPtr = good.ResponseRecPtr

	for i = 0; i < C.MQ_SECURITY_ID_LENGTH; i++ {
		mqod.AlternateSecurityId[i] = C.MQBYTE(good.AlternateSecurityId[i])
	}

	setMQIString((*C.char)(&mqod.ResolvedQName[0]), good.ResolvedQName, C.MQ_OBJECT_NAME_LENGTH)
	setMQIString((*C.char)(&mqod.ResolvedQMgrName[0]), good.ResolvedQMgrName, C.MQ_OBJECT_NAME_LENGTH)

	mqod.ObjectString.VSLength = (C.MQLONG)(len(good.ObjectString))
	mqod.ObjectString.VSCCSID = C.MQCCSI_APPL
	if mqod.ObjectString.VSLength == 0 {
		mqod.ObjectString.VSPtr = C.MQPTR(C.malloc(vsbufsize))
		mqod.ObjectString.VSBufSize = vsbufsize
	} else {
		mqod.ObjectString.VSPtr = (C.MQPTR)(C.CString(good.ObjectString))
	}

	mqod.SelectionString.VSLength = (C.MQLONG)(len(good.SelectionString))
	mqod.SelectionString.VSCCSID = C.MQCCSI_APPL
	if mqod.SelectionString.VSLength == 0 {
		mqod.SelectionString.VSPtr = C.MQPTR(C.malloc(vsbufsize))
		mqod.SelectionString.VSBufSize = vsbufsize
	} else {
		mqod.SelectionString.VSPtr = (C.MQPTR)(C.CString(good.SelectionString))
	}

	mqod.ResObjectString.VSLength = (C.MQLONG)(len(good.ResObjectString))
	mqod.ResObjectString.VSCCSID = C.MQCCSI_APPL
	if mqod.ResObjectString.VSLength == 0 {
		mqod.ResObjectString.VSPtr = C.MQPTR(C.malloc(vsbufsize))
		mqod.ResObjectString.VSBufSize = vsbufsize
	} else {
		mqod.ResObjectString.VSPtr = (C.MQPTR)(C.CString(good.ResObjectString))
	}

	mqod.ResolvedType = good.ResolvedType

	return
}

func copyODfromC(mqod *C.MQOD, good *MQOD) {
	var i int

	good.StrucId = C.GoStringN((*C.char)(&mqod.StrucId[0]), 4)
	good.Version = mqod.Version
	good.ObjectType = mqod.ObjectType
	good.ObjectName = C.GoStringN((*C.char)(&mqod.ObjectName[0]), C.MQ_OBJECT_NAME_LENGTH)
	good.ObjectQMgrName = C.GoStringN((*C.char)(&mqod.ObjectQMgrName[0]), C.MQ_OBJECT_NAME_LENGTH)
	good.DynamicQName = C.GoStringN((*C.char)(&mqod.DynamicQName[0]), C.MQ_OBJECT_NAME_LENGTH)
	good.AlternateUserId = C.GoStringN((*C.char)(&mqod.AlternateUserId[0]), C.MQ_USER_ID_LENGTH)

	good.RecsPresent = mqod.RecsPresent
	good.KnownDestCount = mqod.KnownDestCount
	good.UnknownDestCount = mqod.UnknownDestCount
	good.InvalidDestCount = mqod.InvalidDestCount
	good.ObjectRecOffset = mqod.ObjectRecOffset
	good.ResponseRecOffset = mqod.ResponseRecOffset

	good.ObjectRecPtr = mqod.ObjectRecPtr
	good.ResponseRecPtr = mqod.ResponseRecPtr

	for i = 0; i < C.MQ_SECURITY_ID_LENGTH; i++ {
		good.AlternateSecurityId[i] = (byte)(mqod.AlternateSecurityId[i])
	}

	good.ResolvedQName = C.GoStringN((*C.char)(&mqod.ResolvedQName[0]), C.MQ_OBJECT_NAME_LENGTH)
	good.ResolvedQMgrName = C.GoStringN((*C.char)(&mqod.ResolvedQMgrName[0]), C.MQ_OBJECT_NAME_LENGTH)

	good.ObjectString = C.GoStringN((*C.char)(mqod.ObjectString.VSPtr), (C.int)(mqod.ObjectString.VSLength))
	C.free(unsafe.Pointer(mqod.ObjectString.VSPtr))
	good.SelectionString = C.GoStringN((*C.char)(mqod.SelectionString.VSPtr), (C.int)(mqod.SelectionString.VSLength))
	C.free(unsafe.Pointer(mqod.SelectionString.VSPtr))
	good.ResObjectString = C.GoStringN((*C.char)(mqod.ResObjectString.VSPtr), (C.int)(mqod.ResObjectString.VSLength))
	C.free(unsafe.Pointer(mqod.ResObjectString.VSPtr))
	good.ResolvedType = mqod.ResolvedType

	return
}
