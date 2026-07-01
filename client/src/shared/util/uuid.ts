const hexLookup = Array.from({ length: 256 }, (_, i) => i.toString(16).padStart(2, '0'))

export function uuidv7(): string {
  const bytes = new Uint8Array(16)
  crypto.getRandomValues(bytes)
  const ts = BigInt(Date.now())

  bytes[0] = Number((ts >> 40n) & 0xffn)
  bytes[1] = Number((ts >> 32n) & 0xffn)
  bytes[2] = Number((ts >> 24n) & 0xffn)
  bytes[3] = Number((ts >> 16n) & 0xffn)
  bytes[4] = Number((ts >> 8n) & 0xffn)
  bytes[5] = Number(ts & 0xffn)

  bytes[6] = (bytes[6] & 0x0f) | 0x70
  bytes[8] = (bytes[8] & 0x3f) | 0x80

  return (
    hexLookup[bytes[0]] + hexLookup[bytes[1]] + hexLookup[bytes[2]] + hexLookup[bytes[3]] + '-' +
    hexLookup[bytes[4]] + hexLookup[bytes[5]] + '-' +
    hexLookup[bytes[6]] + hexLookup[bytes[7]] + '-' +
    hexLookup[bytes[8]] + hexLookup[bytes[9]] + '-' +
    hexLookup[bytes[10]] + hexLookup[bytes[11]] + hexLookup[bytes[12]] +
    hexLookup[bytes[13]] + hexLookup[bytes[14]] + hexLookup[bytes[15]]
  )
}
