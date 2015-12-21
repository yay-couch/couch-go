package test

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
    util.Dumpf("Uuid Default       >> %s", uuid.Generate(uuid.HEX_32))
}

/**
 * TestRfc
 */
func TestRfc() {
    util.Dumpf("Uuid RFC           >> %s", uuid.Generate(uuid.RFC))
}

/**
 * TestHex8
 */
func TestHex8() {
    util.Dumpf("Uuid Hex 8         >> %s", uuid.Generate(uuid.HEX_8))
}

/**
 * TestHex32
 */
func TestHex32() {
    util.Dumpf("Uuid Hex 32        >> %s", uuid.Generate(uuid.HEX_32))
}

/**
 * TestHex40
 */
func TestHex40() {
    util.Dumpf("Uuid Hex 40        >> %s", uuid.Generate(uuid.HEX_40))
}

/**
 * TestTimestamp
 */
func TestTimestamp() {
    util.Dumpf("Uuid Timestamp     >> %s", uuid.Generate(uuid.TIMESTAMP))
}

/**
 * TestTimestampNano
 */
func TestTimestampNano() {
    util.Dumpf("Uuid TimestampNano >> %s", uuid.Generate(uuid.TIMESTAMP_NANO))
}
