package test_uuid

import _uuid "./../../src/couch/uuid"

import u "./../../src/couch/util"
// @tmp
var _dump, _dumps, _dumpf = u.Dump, u.Dumps, u.Dumpf

func Shutup() {}

/**
 * TestAll
 */
func TestAll() {
    TestDefault()
    TestRfc()
    TestHex8()
    TestHex32()
    TestHex40()
    TestTimestamp()
    TestTimestampNano()
}

/**
 * TestDefault
 */
func TestDefault() {
    _dumpf("Uuid Default       >> %s", _uuid.GenerateDefault())
}

/**
 * TestRfc
 */
func TestRfc() {
    _dumpf("Uuid RFC           >> %s", _uuid.GenerateRfc())
}

/**
 * TestHex8
 */
func TestHex8() {
    _dumpf("Uuid Hex 8         >> %s", _uuid.GenerateHex8())
}

/**
 * TestHex32
 */
func TestHex32() {
    _dumpf("Uuid Hex 32        >> %s", _uuid.GenerateHex32())
}

/**
 * TestHex40
 */
func TestHex40() {
    _dumpf("Uuid Hex 40        >> %s", _uuid.GenerateHex40())
}

/**
 * TestTimestamp
 */
func TestTimestamp() {
    _dumpf("Uuid Timestamp     >> %s", _uuid.GenerateTimestamp())
}

/**
 * TestTimestampNano
 */
func TestTimestampNano() {
    _dumpf("Uuid TimestampNano >> %s", _uuid.GenerateTimestampNano())
}
