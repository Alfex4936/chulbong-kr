# transform_wgs84_to_korea_tm.bend


# Compute sine using Taylor series with recursion
def sin_taylor(x, term, sum, n, sign):
  if n > 15:
    return sum
  else:
    sum = sum + sign * term
    term = term * x * x / ((2.0 * n) * (2.0 * n + 1.0))
    sign = sign * -1
    return sin_taylor(x, term, sum, n + 1, sign)

def sin(x):
  return sin_taylor(x, x, 0.0, 1, 1.0)

# Compute cosine using Taylor series with recursion
def cos_taylor(x, term, sum, n, sign):
  if n > 15:
    return sum
  else:
    sum = sum + sign * term
    term = term * x * x / ((2.0 * n - 1.0) * (2.0 * n))
    sign = sign * -1
    return cos_taylor(x, term, sum, n + 1, sign)

def cos(x):
  return cos_taylor(x, 1.0, 0.0, 1, 1.0)

# Compute square root using Newton's method with recursion
def sqrt_newton(x, guess):
  new_guess = 0.5 * (guess + x / guess)
  if abs(new_guess - guess) < 0.00001:
    return new_guess
  else:
    return sqrt_newton(x, new_guess)

def sqrt(x):
  return sqrt_newton(x, x / 2.0)

# Compute absolute value
def abs(x):
  if x < 0.0:
    return x * -1
  else:
    return x
def radians(degrees):
  return degrees * 0.0174533  # π / 180

def prepareParams(d, e):
  if e > 1:
    w = e
  else:
    w = 1.0 / e
  z = d * (w - 1) / w
  G = 1.0 - (z * z) / (d * d)
  w = (d * d - z * z) / (z * z)
  z = (d - z) / (d + z)
  return (w, z, G)

def computeE(d, z):
  return d * (1 - z + 5 * (z * z - z * z * z) / 4 + 81 * (z * z * z * z - z * z * z * z * z) / 64)

def computeI(d, z):
  return 3 * d * (z - z * z + 7 * (z * z * z - z * z * z * z) / 8 + 55 * z * z * z * z * z / 64) / 2

def computeJ(d, z):
  return 15 * d * (z * z - z * z * z + 3 * (z * z * z * z - z * z * z * z * z) / 4) / 16

def computeL(d, z):
  return 35 * d * (z * z * z - z * z * z * z + 11 * z * z * z * z * z / 16) / 48

def computeM(d, z):
  return 315 * d * (z * z * z * z - z * z * z * z * z) / 512

def computeIJLM(d, z):
  I = computeI(d, z)
  J = computeJ(d, z)
  L = computeL(d, z)
  M = computeM(d, z)
  return (I, J, L, M)

def computeG(d, G, sinLat):
  return d / sqrt(1 - G * sinLat * sinLat)

def computeU(E, I, J, L, M, latRad):
  return E * latRad - I * sin(2 * latRad) + J * sin(4 * latRad) - L * sin(6 * latRad) + M * sin(8 * latRad)

def computeEComponent(G, sinLat, cosLat, c):
  return G * sinLat * cosLat * c * 0.5

def computeIComponent(G, sinLat, cosLat, t, w, c):
  return G * sinLat * cosLat * cosLat * cosLat * c * (5 - t * t + 9 * w + 4 * w * w) * (1.0 / 24)

def computeJComponent-1(t, w):
  return 61 - 58 * t * t + t * t * t * t + 270 * w - 330 * t * t * w + 445 * w * w + 324 * w * w * w - computeJComponent-2(t, w)
def computeJComponent-2(t, w):
  return 680 * t * t * w * w + 88 * w * w * w * w - 600 * t * t * w * w * w - 192 * t * t * w * w * w * w

def computeJComponent(G, sinLat, cosLat, t, w, c):
  return G * sinLat * cosLat * cosLat * cosLat * cosLat * c * (computeJComponent-1(t,w)) * (1.0 / 720)

def computeHComponent-1(t):
  return 1385 - 3111 * t * t + 543 * t * t * t * t - t * t * t * t * t * t
def computeHComponent(G, sinLat, cosLat, t, c):
  return G * sinLat * cosLat * cosLat * cosLat * cosLat * cosLat * cosLat * c * (computeHComponent-1(t)) * (1.0 / 40320)

def computeO-1(u, c, D, E, I, J, H):
  return u * c + D * D * E + D * D * D * I + D * D * D * D * D * J + computeO-2(D, H)

def computeO-2(D, H):
  return D * D * D * D * D * D * D * H

def computeO(D, G, sinLat, cosLat, t, w, u, c):
  E = computeEComponent(G, sinLat, cosLat, c)
  I = computeIComponent(G, sinLat, cosLat, t, w, c)
  J = computeJComponent(G, sinLat, cosLat, t, w, c)
  H = computeHComponent(G, sinLat, cosLat, t, c)
  o = computeO-1(u, c, D, E, I, J, H)
  return o

def computeYComponent(o, h):
  return o - h

def computeXComponent1(D, G, cosLat, c):
  return G * cosLat * c

def computeXComponent2(D, G, cosLat, t, w, c):
  return G * cosLat * cosLat * cosLat * c * (1 - t * t + w) * (1.0 / 6)


def computeXComponent3-1(t, w):
  return 5 - 18 * t * t + t * t * t * t + 14 * w - 58 * t * t * w + 13 * w * w + 4 * w * w * w
def computeXComponent3-2(t, w):
  return 64 * t * t * w * w - 25 * t * t * w * w * w

def computeXComponent3(D, G, cosLat, t, w, c):
  return G * cosLat * cosLat * cosLat * cosLat * cosLat * c * (computeXComponent3-1(t, w) - computeXComponent3-2(t, w)) * (1.0 / 120)

def computeXComponent4-1(t):
  return 61 - 479 * t * t + 179 * t * t * t * t - t * t * t * t * t * t

def computeXComponent4(D, G, cosLat, t, c):
  return G * cosLat * cosLat * cosLat * cosLat * cosLat * cosLat * cosLat * c * computeXComponent4-1(t) * (1.0 / 5040)

def computeX(f, D, E, z, w, u):
  return f + D * E + D * D * D * z + D * D * D * D * D * w + D * D * D * D * D * D * D * u

def computeXY(D, G, cosLat, t, w, c, o, f, h):
  y = computeYComponent(o, h)
  E = computeXComponent1(D, G, cosLat, c)
  z = computeXComponent2(D, G, cosLat, t, w, c)
  w = computeXComponent3(D, G, cosLat, t, w, c)
  u = computeXComponent4(D, G, cosLat, t, c)
  x = computeX(f, D, E, z, w, u)
  return (x, y)

def getRadians(lat, lon, l, m):
  return (radians(lat), radians(lon), radians(l), radians(m))
  
def transformWGS84ToKoreaTM(d, e, h, f, c, l, m, lat, lon):
  (latRad, lonRad, lRad, mRad) = getRadians(lat, lon, l, m)
  (w, z, G) = prepareParams(d, e)
  E = computeE(d, z)
  (I, J, L, M) = computeIJLM(d, z)
  D = lonRad - mRad
  (sinLat, cosLat) = computeSinCos(latRad)
  t = sinLat / cosLat
  G = computeG(d, G, sinLat)
  u = computeU(E, I, J, L, M, latRad)
  o = computeO(D, G, sinLat, cosLat, t, w, u, c)
  return computeXY(D, G, cosLat, t, w, c, o, f, h)

def computeSinCos(latRad):
  sinLat = sin(latRad)
  cosLat = cos(latRad)
  return (sinLat, cosLat)

def main:
  d = 6378137.0
  e = 1 / 298.257223563
  h = 200000.0
  f = 500000.0
  c = 1.0
  l = 38.0
  m = 127.5
  lat = 37.5665
  lon = 126.9780

  (x, y) = transformWGS84ToKoreaTM(d, e, h, f, c, l, m, lat, lon)
  return (x, y)