// File "very_simple_switch_model.p4"
// Very Simple Switch P4 declaration

#include <core.p4>

/* Various constants and structure declarations */
typedef bit<4> PortId;  // Ports are represented using 4-bit values
const PortId REAL_PORT_COUNT = 4w8;  // Number of real ports (8)

/* Metadata accompanying an input packet */
struct InControl {
    PortId inputPort;
};

/* Special input port values */
const PortId RECIRCULATE_IN_PORT = 0xD;
const PortId CPU_IN_PORT = 0xE;

/* Metadata that must be computed for outgoing packets */
struct OutControl {
    PortId outputPort;
};

/* Special output port values for outgoing packet */
const PortId DROP_PORT = 0xF;
const PortId CPU_OUT_PORT = 0xE;
const PortId RECIRCULATE_OUT_PORT = 0xD;

/* Prototypes for all programmable blocks */

/**
 * Programmable parser.
 * @param <H> type of headers; defined by user
 * @param b input packet
 * @param parsedHeaders headers constructed by parser
 */
parser Parser<H>(packet_in b,
                 out H parsedHeaders);

/**
 * Match-action pipeline
 * @param <H> type of input and output headers
 * @param headers headers received from the parser and sent to the deparser
 * @param parseError error that may have surfaced during parsing
 * @param inCtrl information from architecture, accompanying input packet
 * @param outCtrl information for architecture, accompanying output packet
 */
control Pipe<H>(inout H headers,
                in error parseError,
                in InControl inCtrl,
                out OutControl outCtrl);

/**
 * VSS deparser.
 * @param <H> type of headers; defined by user
 * @param b output packet
 * @param outputHeaders headers for output packet
 */
control Deparser<H>(inout H outputHeaders,
                    packet_out b);

/**
 * Top-level package declaration - must be instantiated by user.
 * @param <H> user-defined type of the headers processed.
 */
package VSS<H>(Parser<H> p,
               Pipe<H> map,
               Deparser<H> d);

// Architecture-specific objects that can be instantiated

// Checksum unit
extern Checksum16 {
    Checksum16();  // Constructor
    void clear();  // Prepare unit for computation
    void update<T>(in T data);  // Add data to checksum
    void remove<T>(in T data);  // Remove data from existing checksum
    bit<16> get();  // Get the checksum for the data added since last clear
}