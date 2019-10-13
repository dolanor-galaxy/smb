package smbnego

import (
	"encoding/binary"
	"time"

	"github.com/gentlemanautomaton/smb/smbcap"
	"github.com/gentlemanautomaton/smb/smbdialect"
	"github.com/gentlemanautomaton/smb/smbid"
	"github.com/gentlemanautomaton/smb/smbsecmode"
	"github.com/gentlemanautomaton/smb/smbtype"
)

// Response interprets a slice of bytes as an SMB negotiation response packet.
//
// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-smb2/63abf97c-0d09-47e2-88d6-6bfa552949a5
type Response []byte

// Valid returns true if the response is valid.
func (r Response) Valid() bool {
	if len(r) < 64 {
		return false
	}

	// The spec requires the size field to be 65
	if r.Size() != 65 {
		return false
	}

	// The security buffer must not overflow
	if int(r.SecurityBufferOffset())+int(r.SecurityBufferLength()) > len(r) {
		return false
	}

	// In SMB 3.1.1 the negotiation contexts must not overflow
	if r.DialectRevision() == smbdialect.SMB311 {
		// Make sure the context count is compatible with the size of the
		// response. The size of each context is variable but at least 8 bytes.
		minimumLength := uint(r.ContextOffset()) + uint(r.ContextCount())*ContextHeaderLength
		if minimumLength > uint(len(r)) {
			return false
		}

		// Rely on the context list implementation to determine its own validity
		if !r.ContextList().Valid(r.ContextCount()) {
			return false
		}
	}

	return true
}

// Size returns the structure size of the response. The specification requires
// that this be 65, regardless of the size of the security buffer or the
// number of negotiation contexts.
func (r Response) Size() uint16 {
	return binary.LittleEndian.Uint16(r[0:2])
}

// SetSize sets the structure size of the response.
func (r Response) SetSize(size uint16) {
	binary.LittleEndian.PutUint16(r[0:2], size)
}

// SecurityMode returns the security mode of the response.
func (r Response) SecurityMode() smbsecmode.Flags {
	return smbsecmode.Flags(binary.LittleEndian.Uint16(r[2:4]))
}

// SetSecurityMode sets the security mode of the response.
func (r Response) SetSecurityMode(flags smbsecmode.Flags) {
	binary.LittleEndian.PutUint16(r[2:4], uint16(flags))
}

// DialectRevision returns the dialect revision of the response.
func (r Response) DialectRevision() smbdialect.Revision {
	return smbdialect.Revision(binary.LittleEndian.Uint16(r[4:6]))
}

// SetDialectRevision sets the dialect revision of the response.
func (r Response) SetDialectRevision(revision smbdialect.Revision) {
	binary.LittleEndian.PutUint16(r[4:6], uint16(revision))
}

// ContextCount returns the context count of the response.
//
// This field is only valid in the SMB 3.1.1 dialect.
func (r Response) ContextCount() uint16 {
	return binary.LittleEndian.Uint16(r[6:8])
}

// SetContextCount sets the context count of the response.
//
// This field is only valid in the SMB 3.1.1 dialect.
func (r Response) SetContextCount(size uint16) {
	binary.LittleEndian.PutUint16(r[6:8], size)
}

// ServerID returns the server identifier of the response.
func (r Response) ServerID() (id smbid.ID) {
	id.Read(r[8:24])
	return
}

// SetServerID sets the server identifier of the response.
func (r Response) SetServerID(id smbid.ID) {
	id.Write(r[8:24])
}

// Capabilities returns the capability flags of the response.
func (r Response) Capabilities() smbcap.Flags {
	return smbcap.Flags(binary.LittleEndian.Uint32(r[24:28]))
}

// SetCapabilities sets the capability flags of the response.
func (r Response) SetCapabilities(flags smbcap.Flags) {
	binary.LittleEndian.PutUint32(r[24:28], uint32(flags))
}

// MaxTransactSize returns the maximum transaction size of the response.
func (r Response) MaxTransactSize() uint32 {
	return binary.LittleEndian.Uint32(r[28:32])
}

// SetMaxTransactSize sets the maximum transaction size of the response.
func (r Response) SetMaxTransactSize(flags uint32) {
	binary.LittleEndian.PutUint32(r[28:32], flags)
}

// MaxReadSize returns the maximum read size of the response.
func (r Response) MaxReadSize() uint32 {
	return binary.LittleEndian.Uint32(r[32:36])
}

// SetMaxReadSize sets the maximum read size of the response.
func (r Response) SetMaxReadSize(flags uint32) {
	binary.LittleEndian.PutUint32(r[32:36], flags)
}

// MaxWriteSize returns the maximum write size of the response.
func (r Response) MaxWriteSize() uint32 {
	return binary.LittleEndian.Uint32(r[36:40])
}

// SetMaxWriteSize sets the maximum write size of the response.
func (r Response) SetMaxWriteSize(flags uint32) {
	binary.LittleEndian.PutUint32(r[36:40], flags)
}

// SystemTime returns the system time of the response.
func (r Response) SystemTime() time.Time {
	return smbtype.Time(r[40:48])
}

// SetSystemTime sets the system time of the response.
func (r Response) SetSystemTime(t time.Time) {
	smbtype.PutTime(r[40:48], t)
}

// ServerStartTime returns the server start time of the response.
func (r Response) ServerStartTime() time.Time {
	return smbtype.Time(r[48:56])
}

// SetServerStartTime sets the server start time of the response.
func (r Response) SetServerStartTime(t time.Time) {
	smbtype.PutTime(r[48:56], t)
}

// SecurityBufferOffset returns the offset of the security buffer within the
// response.
func (r Response) SecurityBufferOffset() uint16 {
	return binary.LittleEndian.Uint16(r[56:58])
}

// SetSecurityBufferOffset sets the offset of the security buffer within the
// response.
func (r Response) SetSecurityBufferOffset(offset uint16) {
	binary.LittleEndian.PutUint16(r[56:58], offset)
}

// SecurityBufferLength returns the length of the security buffer within the
// response.
func (r Response) SecurityBufferLength() uint16 {
	return binary.LittleEndian.Uint16(r[58:60])
}

// SetSecurityBufferLength sets the length of the security buffer within the
// response.
func (r Response) SetSecurityBufferLength(length uint16) {
	binary.LittleEndian.PutUint16(r[58:60], length)
}

// SecurityBuffer returns the bytes of the security buffer from the response.
func (r Response) SecurityBuffer() []byte {
	start := uint(r.SecurityBufferOffset())
	length := uint(r.SecurityBufferLength())
	end := start + length
	return r[start:end:end]
}

// ContextOffset returns the offset of the first negotiate context
// within the response.
//
// This field is only valid in the SMB 3.1.1 dialect.
func (r Response) ContextOffset() uint32 {
	return binary.LittleEndian.Uint32(r[60:64])
}

// SetContextOffset sets the offset of the first negotiate context
// within the response.
//
// This field is only valid in the SMB 3.1.1 dialect.
func (r Response) SetContextOffset(size uint32) {
	binary.LittleEndian.PutUint32(r[60:64], size)
}

// ContextList returns the negotiation context list from the response.
//
// If r is valid the returned list is guaranteed to be valid.
//
// This field is only valid in the SMB 3.1.1 dialect.
func (r Response) ContextList() ContextList {
	return ContextList(r[r.ContextOffset():])
}
