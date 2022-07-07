# import matplotlib.pyplot as plt
import pandas as pd

crime = pd.read_csv('crime.csv')
police = pd.read_csv('police.csv')
police['country'] = \
    police['Country or dependency'].str.replace(r'\s+\*.*', '', regex=True)

# TODO: More country normalization
df = pd.merge(crime, police, left_on='country', right_on='country')
xs = 'police per 100k'
df = df.rename(columns={
    'Rate.mw-parser-output .nobold{font-weight:normal}(per 100k people)': xs,
})
df[xs] = df[xs].str.replace(',', '').astype('int64')
# df.plot(x=xs, y='crimeIndex')
# plt.show()
corr = df[xs].corr(df['crimeIndex'])
print(f'correlation between % of police to crime index: {corr:.2f}')
