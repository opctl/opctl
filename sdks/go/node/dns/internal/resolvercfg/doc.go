// package resolvercfg allows configuring DNS resolution for a given OS.
//
// some OS's, like OSX, support hierarchical resolver configs whereby each domain can
// have It's own config (see man 5 resolver); Others, like linux, only support a single global resolver config.

// In order to shield consumers from these inconsistencies, we take a least common denominator
// approach to our exposed interface i.e. expose per domain resolver config.
package resolvercfg
