Results of running
```
go test -count=5 -bench=././././serial | benchstat -html > results.md
```

The benchmark names follow the format `Adapter/(body_size)/(encoder)/(level)/serial-*`.

The first table lists time to compress the body with the specified encoder and level.

The second table lists the size of the resulting compressed body.

<style>.benchstat tbody td:nth-child(1n+2) { text-align: right; padding: 0em 1em; }</style>
<table class='benchstat'>
<tr><th>name</th><th>time/op</th>
<tr><td>Adapter/10/gzip/1/serial-12</td><td>924ns ± 0%</td>
<tr><td>Adapter/10/br/1/serial-12</td><td>933ns ± 5%</td>
<tr><td>Adapter/10/zstd/1/serial-12</td><td>968ns ± 5%</td>
<tr><td>Adapter/100/gzip/1/serial-12</td><td>960ns ± 6%</td>
<tr><td>Adapter/100/br/1/serial-12</td><td>963ns ± 6%</td>
<tr><td>Adapter/100/zstd/1/serial-12</td><td>975ns ± 6%</td>
<tr><td>Adapter/1000/gzip/1/serial-12</td><td>30.8µs ±10%</td>
<tr><td>Adapter/1000/gzip/2/serial-12</td><td>85.9µs ± 5%</td>
<tr><td>Adapter/1000/gzip/3/serial-12</td><td>85.6µs ± 0%</td>
<tr><td>Adapter/1000/gzip/4/serial-12</td><td>89.2µs ± 5%</td>
<tr><td>Adapter/1000/gzip/5/serial-12</td><td>89.3µs ± 0%</td>
<tr><td>Adapter/1000/gzip/6/serial-12</td><td>91.4µs ± 2%</td>
<tr><td>Adapter/1000/gzip/7/serial-12</td><td>93.1µs ± 6%</td>
<tr><td>Adapter/1000/gzip/8/serial-12</td><td>93.4µs ± 6%</td>
<tr><td>Adapter/1000/gzip/9/serial-12</td><td>93.5µs ± 8%</td>
<tr><td>Adapter/1000/br/1/serial-12</td><td>27.7µs ± 1%</td>
<tr><td>Adapter/1000/br/2/serial-12</td><td>42.6µs ±10%</td>
<tr><td>Adapter/1000/br/3/serial-12</td><td>69.9µs ± 4%</td>
<tr><td>Adapter/1000/br/4/serial-12</td><td>99.4µs ±12%</td>
<tr><td>Adapter/1000/br/5/serial-12</td><td>111µs ± 7%</td>
<tr><td>Adapter/1000/br/6/serial-12</td><td>112µs ± 4%</td>
<tr><td>Adapter/1000/br/7/serial-12</td><td>141µs ± 4%</td>
<tr><td>Adapter/1000/br/8/serial-12</td><td>147µs ±19%</td>
<tr><td>Adapter/1000/br/9/serial-12</td><td>154µs ± 1%</td>
<tr><td>Adapter/1000/br/10/serial-12</td><td>2.44ms ± 5%</td>
<tr><td>Adapter/1000/br/11/serial-12</td><td>4.09ms ± 1%</td>
<tr><td>Adapter/1000/zstd/1/serial-12</td><td>18.1µs ±10%</td>
<tr><td>Adapter/1000/zstd/2/serial-12</td><td>21.6µs ± 7%</td>
<tr><td>Adapter/1000/zstd/3/serial-12</td><td>31.2µs ±14%</td>
<tr><td>Adapter/1000/zstd/4/serial-12</td><td>79.9µs ±14%</td>
<tr><td>Adapter/10000/gzip/1/serial-12</td><td>134µs ±19%</td>
<tr><td>Adapter/10000/gzip/2/serial-12</td><td>201µs ± 3%</td>
<tr><td>Adapter/10000/gzip/3/serial-12</td><td>203µs ± 1%</td>
<tr><td>Adapter/10000/gzip/4/serial-12</td><td>247µs ±19%</td>
<tr><td>Adapter/10000/gzip/5/serial-12</td><td>246µs ± 1%</td>
<tr><td>Adapter/10000/gzip/6/serial-12</td><td>269µs ± 7%</td>
<tr><td>Adapter/10000/gzip/7/serial-12</td><td>335µs ± 3%</td>
<tr><td>Adapter/10000/gzip/8/serial-12</td><td>371µs ± 0%</td>
<tr><td>Adapter/10000/gzip/9/serial-12</td><td>381µs ±11%</td>
<tr><td>Adapter/10000/br/1/serial-12</td><td>148µs ± 2%</td>
<tr><td>Adapter/10000/br/2/serial-12</td><td>325µs ±10%</td>
<tr><td>Adapter/10000/br/3/serial-12</td><td>389µs ± 0%</td>
<tr><td>Adapter/10000/br/4/serial-12</td><td>557µs ±19%</td>
<tr><td>Adapter/10000/br/5/serial-12</td><td>614µs ± 2%</td>
<tr><td>Adapter/10000/br/6/serial-12</td><td>655µs ± 8%</td>
<tr><td>Adapter/10000/br/7/serial-12</td><td>804µs ± 4%</td>
<tr><td>Adapter/10000/br/8/serial-12</td><td>823µs ± 7%</td>
<tr><td>Adapter/10000/br/9/serial-12</td><td>978µs ± 5%</td>
<tr><td>Adapter/10000/br/10/serial-12</td><td>10.9ms ± 8%</td>
<tr><td>Adapter/10000/br/11/serial-12</td><td>29.5ms ±14%</td>
<tr><td>Adapter/10000/zstd/1/serial-12</td><td>74.6µs ±16%</td>
<tr><td>Adapter/10000/zstd/2/serial-12</td><td>124µs ± 8%</td>
<tr><td>Adapter/10000/zstd/3/serial-12</td><td>200µs ± 7%</td>
<tr><td>Adapter/10000/zstd/4/serial-12</td><td>634µs ± 8%</td>
<tr><td>Adapter/100000/gzip/1/serial-12</td><td>1.12ms ± 9%</td>
<tr><td>Adapter/100000/gzip/2/serial-12</td><td>1.35ms ± 1%</td>
<tr><td>Adapter/100000/gzip/3/serial-12</td><td>1.53ms ± 5%</td>
<tr><td>Adapter/100000/gzip/4/serial-12</td><td>1.73ms ± 6%</td>
<tr><td>Adapter/100000/gzip/5/serial-12</td><td>1.99ms ± 1%</td>
<tr><td>Adapter/100000/gzip/6/serial-12</td><td>2.34ms ± 7%</td>
<tr><td>Adapter/100000/gzip/7/serial-12</td><td>3.80ms ±12%</td>
<tr><td>Adapter/100000/gzip/8/serial-12</td><td>8.67ms ± 3%</td>
<tr><td>Adapter/100000/gzip/9/serial-12</td><td>8.96ms ±12%</td>
<tr><td>Adapter/100000/br/1/serial-12</td><td>1.48ms ±13%</td>
<tr><td>Adapter/100000/br/2/serial-12</td><td>3.01ms ± 5%</td>
<tr><td>Adapter/100000/br/3/serial-12</td><td>3.51ms ± 2%</td>
<tr><td>Adapter/100000/br/4/serial-12</td><td>4.58ms ± 7%</td>
<tr><td>Adapter/100000/br/5/serial-12</td><td>5.51ms ± 0%</td>
<tr><td>Adapter/100000/br/6/serial-12</td><td>6.04ms ± 4%</td>
<tr><td>Adapter/100000/br/7/serial-12</td><td>7.11ms ± 0%</td>
<tr><td>Adapter/100000/br/8/serial-12</td><td>7.86ms ± 3%</td>
<tr><td>Adapter/100000/br/9/serial-12</td><td>10.2ms ± 5%</td>
<tr><td>Adapter/100000/br/10/serial-12</td><td>113ms ±10%</td>
<tr><td>Adapter/100000/br/11/serial-12</td><td>312ms ± 2%</td>
<tr><td>Adapter/100000/zstd/1/serial-12</td><td>629µs ± 1%</td>
<tr><td>Adapter/100000/zstd/2/serial-12</td><td>859µs ± 4%</td>
<tr><td>Adapter/100000/zstd/3/serial-12</td><td>1.02ms ± 2%</td>
<tr><td>Adapter/100000/zstd/4/serial-12</td><td>3.79ms ±11%</td>
</table>

<style>.benchstat tbody td:nth-child(1n+2) { text-align: right; padding: 0em 1em; }</style>
<table class='benchstat'>
<tr><th>name</th><th>%</th>
<tr><td>Adapter/10/gzip/1/serial-12</td><td>100 ± 0%</td>
<tr><td>Adapter/10/br/1/serial-12</td><td>100 ± 0%</td>
<tr><td>Adapter/10/zstd/1/serial-12</td><td>100 ± 0%</td>
<tr><td>Adapter/100/gzip/1/serial-12</td><td>100 ± 0%</td>
<tr><td>Adapter/100/br/1/serial-12</td><td>100 ± 0%</td>
<tr><td>Adapter/100/zstd/1/serial-12</td><td>100 ± 0%</td>
<tr><td>Adapter/1000/gzip/1/serial-12</td><td>44.2 ± 0%</td>
<tr><td>Adapter/1000/gzip/2/serial-12</td><td>41.9 ± 0%</td>
<tr><td>Adapter/1000/gzip/3/serial-12</td><td>41.7 ± 0%</td>
<tr><td>Adapter/1000/gzip/4/serial-12</td><td>41.7 ± 0%</td>
<tr><td>Adapter/1000/gzip/5/serial-12</td><td>41.1 ± 0%</td>
<tr><td>Adapter/1000/gzip/6/serial-12</td><td>41.1 ± 0%</td>
<tr><td>Adapter/1000/gzip/7/serial-12</td><td>41.0 ± 0%</td>
<tr><td>Adapter/1000/gzip/8/serial-12</td><td>41.0 ± 0%</td>
<tr><td>Adapter/1000/gzip/9/serial-12</td><td>41.0 ± 0%</td>
<tr><td>Adapter/1000/br/1/serial-12</td><td>45.2 ± 0%</td>
<tr><td>Adapter/1000/br/2/serial-12</td><td>42.1 ± 0%</td>
<tr><td>Adapter/1000/br/3/serial-12</td><td>39.9 ± 0%</td>
<tr><td>Adapter/1000/br/4/serial-12</td><td>39.3 ± 0%</td>
<tr><td>Adapter/1000/br/5/serial-12</td><td>36.6 ± 0%</td>
<tr><td>Adapter/1000/br/6/serial-12</td><td>36.8 ± 0%</td>
<tr><td>Adapter/1000/br/7/serial-12</td><td>36.7 ± 0%</td>
<tr><td>Adapter/1000/br/8/serial-12</td><td>36.7 ± 0%</td>
<tr><td>Adapter/1000/br/9/serial-12</td><td>36.7 ± 0%</td>
<tr><td>Adapter/1000/br/10/serial-12</td><td>37.4 ± 0%</td>
<tr><td>Adapter/1000/br/11/serial-12</td><td>37.2 ± 0%</td>
<tr><td>Adapter/1000/zstd/1/serial-12</td><td>43.1 ± 0%</td>
<tr><td>Adapter/1000/zstd/2/serial-12</td><td>42.2 ± 0%</td>
<tr><td>Adapter/1000/zstd/3/serial-12</td><td>41.7 ± 0%</td>
<tr><td>Adapter/1000/zstd/4/serial-12</td><td>41.3 ± 0%</td>
<tr><td>Adapter/10000/gzip/1/serial-12</td><td>29.6 ± 0%</td>
<tr><td>Adapter/10000/gzip/2/serial-12</td><td>28.5 ± 0%</td>
<tr><td>Adapter/10000/gzip/3/serial-12</td><td>28.1 ± 0%</td>
<tr><td>Adapter/10000/gzip/4/serial-12</td><td>28.2 ± 0%</td>
<tr><td>Adapter/10000/gzip/5/serial-12</td><td>27.7 ± 0%</td>
<tr><td>Adapter/10000/gzip/6/serial-12</td><td>27.6 ± 0%</td>
<tr><td>Adapter/10000/gzip/7/serial-12</td><td>27.2 ± 0%</td>
<tr><td>Adapter/10000/gzip/8/serial-12</td><td>27.2 ± 0%</td>
<tr><td>Adapter/10000/gzip/9/serial-12</td><td>27.2 ± 0%</td>
<tr><td>Adapter/10000/br/1/serial-12</td><td>29.9 ± 0%</td>
<tr><td>Adapter/10000/br/2/serial-12</td><td>28.1 ± 0%</td>
<tr><td>Adapter/10000/br/3/serial-12</td><td>27.9 ± 0%</td>
<tr><td>Adapter/10000/br/4/serial-12</td><td>27.5 ± 0%</td>
<tr><td>Adapter/10000/br/5/serial-12</td><td>26.0 ± 0%</td>
<tr><td>Adapter/10000/br/6/serial-12</td><td>25.9 ± 0%</td>
<tr><td>Adapter/10000/br/7/serial-12</td><td>25.8 ± 0%</td>
<tr><td>Adapter/10000/br/8/serial-12</td><td>25.8 ± 0%</td>
<tr><td>Adapter/10000/br/9/serial-12</td><td>25.8 ± 0%</td>
<tr><td>Adapter/10000/br/10/serial-12</td><td>23.4 ± 0%</td>
<tr><td>Adapter/10000/br/11/serial-12</td><td>23.1 ± 0%</td>
<tr><td>Adapter/10000/zstd/1/serial-12</td><td>28.6 ± 0%</td>
<tr><td>Adapter/10000/zstd/2/serial-12</td><td>28.1 ± 0%</td>
<tr><td>Adapter/10000/zstd/3/serial-12</td><td>28.0 ± 0%</td>
<tr><td>Adapter/10000/zstd/4/serial-12</td><td>27.6 ± 0%</td>
<tr><td>Adapter/100000/gzip/1/serial-12</td><td>27.3 ± 0%</td>
<tr><td>Adapter/100000/gzip/2/serial-12</td><td>26.2 ± 0%</td>
<tr><td>Adapter/100000/gzip/3/serial-12</td><td>25.7 ± 0%</td>
<tr><td>Adapter/100000/gzip/4/serial-12</td><td>25.8 ± 0%</td>
<tr><td>Adapter/100000/gzip/5/serial-12</td><td>25.0 ± 0%</td>
<tr><td>Adapter/100000/gzip/6/serial-12</td><td>24.8 ± 0%</td>
<tr><td>Adapter/100000/gzip/7/serial-12</td><td>24.2 ± 0%</td>
<tr><td>Adapter/100000/gzip/8/serial-12</td><td>24.1 ± 0%</td>
<tr><td>Adapter/100000/gzip/9/serial-12</td><td>24.1 ± 0%</td>
<tr><td>Adapter/100000/br/1/serial-12</td><td>27.6 ± 0%</td>
<tr><td>Adapter/100000/br/2/serial-12</td><td>25.4 ± 0%</td>
<tr><td>Adapter/100000/br/3/serial-12</td><td>25.3 ± 0%</td>
<tr><td>Adapter/100000/br/4/serial-12</td><td>24.9 ± 0%</td>
<tr><td>Adapter/100000/br/5/serial-12</td><td>23.6 ± 0%</td>
<tr><td>Adapter/100000/br/6/serial-12</td><td>23.4 ± 0%</td>
<tr><td>Adapter/100000/br/7/serial-12</td><td>23.3 ± 0%</td>
<tr><td>Adapter/100000/br/8/serial-12</td><td>23.2 ± 0%</td>
<tr><td>Adapter/100000/br/9/serial-12</td><td>23.2 ± 0%</td>
<tr><td>Adapter/100000/br/10/serial-12</td><td>20.2 ± 0%</td>
<tr><td>Adapter/100000/br/11/serial-12</td><td>19.9 ± 0%</td>
<tr><td>Adapter/100000/zstd/1/serial-12</td><td>25.9 ± 0%</td>
<tr><td>Adapter/100000/zstd/2/serial-12</td><td>25.7 ± 0%</td>
<tr><td>Adapter/100000/zstd/3/serial-12</td><td>25.6 ± 0%</td>
<tr><td>Adapter/100000/zstd/4/serial-12</td><td>25.3 ± 0%</td>
</table>
