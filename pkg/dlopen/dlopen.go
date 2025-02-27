package dlopen

// #cgo LDFLAGS: -ldl
// #include <stdlib.h>
// #include <dlfcn.h>
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

// LibHandle represents an open handle to a library (.so)
type LibHandle struct {
	Handle  unsafe.Pointer
	Libname string
}

// GetHandle tries to get a handle to a library (.so), attempting to access it
// by the names specified in libs and returning the first that is successfully
// opened. Callers are responsible for closing the handler. If no library can
// be successfully opened, an error is returned.
func GetHandle(libs []string) (*LibHandle, error) {
	for _, name := range libs {
		libname := C.CString(name)
		defer C.free(unsafe.Pointer(libname))
		handle := C.dlopen(libname, C.RTLD_LAZY)
		if handle != nil {
			h := &LibHandle{
				Handle:  handle,
				Libname: name,
			}
			return h, nil
		}
	}

	return nil, fmt.Errorf(C.GoString(C.dlerror()))
}

// GetSymbolPointer takes a symbol name and returns a pointer to the symbol.
func (l *LibHandle) GetSymbolPointer(symbol string) (unsafe.Pointer, error) {
	sym := C.CString(symbol)
	defer C.free(unsafe.Pointer(sym))

	C.dlerror()
	p := C.dlsym(l.Handle, sym)
	e := C.dlerror()
	if e != nil {
		return nil, fmt.Errorf("error resolving symbol %q: %v", symbol, errors.New(C.GoString(e)))
	}

	return p, nil
}

// Close closes a LibHandle.
func (l *LibHandle) Close() error {
	C.dlerror()
	C.dlclose(l.Handle)
	e := C.dlerror()
	if e != nil {
		return fmt.Errorf("error closing %v: %v", l.Libname, errors.New(C.GoString(e)))
	}

	return nil
}
