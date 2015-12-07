package uuid

import (
    "./../../src/couch/util"
    "./../../src/couch/uuid"
)

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
    util.Dumpf("Uuid Default       >> %s", uuid.GenerateDefault())
}

/**
 * TestRfc
 */
func TestRfc() {
    util.Dumpf("Uuid RFC           >> %s", uuid.GenerateRfc())
}

/**
 * TestHex8
 */
func TestHex8() {
    util.Dumpf("Uuid Hex 8         >> %s", uuid.GenerateHex8())
}

/**
 * TestHex32
 */
func TestHex32() {
    util.Dumpf("Uuid Hex 32        >> %s", uuid.GenerateHex32())
}

/**
 * TestHex40
 */
func TestHex40() {
    util.Dumpf("Uuid Hex 40        >> %s", uuid.GenerateHex40())
}

/**
 * TestTimestamp
 */
func TestTimestamp() {
    util.Dumpf("Uuid Timestamp     >> %s", uuid.GenerateTimestamp())
}

/**
 * TestTimestampNano
 */
func TestTimestampNano() {
    util.Dumpf("Uuid TimestampNano >> %s", uuid.GenerateTimestampNano())
}
