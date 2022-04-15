from scipy.optimize import curve_fit
import numpy as np
import matplotlib.pyplot as plt

data = np.genfromtxt('data.txt',delimiter=',')
x, y = data.T

# plt.plot(x,y)
# plt.show()

def func(x, *params):
    y = np.zeros_like(x)
    for i in range(0, len(params), 3):
        ctr = params[i]
        amp = params[i+1]
        wid = params[i+2]
        y = y + amp * np.exp( -((x - ctr)/wid)**2)
    return y

guess = [110,50000,10]
guess += [269,200,10]

popt, pcov = curve_fit(func, x, y, p0=guess)
fit = func(x, *popt)
print(popt)
plt.plot(x, y)
plt.plot(x, fit , 'r-')
plt.show()