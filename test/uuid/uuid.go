package uuid

import (
    "./../../src/couch/util"
    "./../../src/couch/uuid"
)

import _uuid "./../../src/couch/uuid"

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
    util.Dumpf("Uuid Default       >> %s", _uuid.GenerateDefault())
}

/**
 * TestRfc
 */
func TestRfc() {
    util.Dumpf("Uuid RFC           >> %s", _uuid.GenerateRfc())
}

/**
 * TestHex8
 */
func TestHex8() {
    util.Dumpf("Uuid Hex 8         >> %s", _uuid.GenerateHex8())
}

/**
 * TestHex32
 */
func TestHex32() {
    util.Dumpf("Uuid Hex 32        >> %s", _uuid.GenerateHex32())
}

/**
 * TestHex40
 */
func TestHex40() {
    util.Dumpf("Uuid Hex 40        >> %s", _uuid.GenerateHex40())
}

/**
 * TestTimestamp
 */
func TestTimestamp() {
    util.Dumpf("Uuid Timestamp     >> %s", _uuid.GenerateTimestamp())
}

/**
 * TestTimestampNano
 */
func TestTimestampNano() {
    util.Dumpf("Uuid TimestampNano >> %s", _uuid.GenerateTimestampNano())
}
