# M3U8 PTS Analyser


The main point of this program is to check the stitched segments from Google DAI and evaluate its TS time metadatas.


Output example:


```
> go run main.go https://dai.google.com/linear/hls/pa/event/IsQ-v81PScyXGohAP6aqLA/stream/1ff053fa-e421-488d-90de-e545f4ce3a8b:SCL/variant/f221d90f556280632ac5fa7c25e081a7/bandwidth/764000.m3u8

Analyzing: https://dai.google.com/linear/hls/pa/event/IsQ-v81PScyXGohAP6aqLA/stream/1ff053fa-e421-488d-90de-e545f4ce3a8b:SCL/variant/f221d90f556280632ac5fa7c25e081a7/bandwidth/764000.m3u8

Saving files at: /tmp/1725678607

DISCONTINUITY
0 | 0.000000 | 450450 | > startTime: 0.000000s endTime: 5.013333s [0.ts]
450450 | 5.005000 | 450450 | > startTime: 5.005000s endTime: 10.010000s [1.ts]
900900 | 10.010000 | 450450 | > startTime: 10.005333s endTime: 15.018666s [2.ts]
1351350 | 15.015000 | 33033 | > startTime: 15.015000s endTime: 15.382033s [3.ts]
DISCONTINUITY
0 | 0.000000 | 450450 | > startTime: 0.000000s endTime: 5.013333s [0.ts]
450450 | 5.005000 | 450450 | > startTime: 5.005000s endTime: 10.010000s [1.ts]
900900 | 10.010000 | 426426 | > startTime: 10.005333s endTime: 14.748067s [2.ts]
DISCONTINUITY
```
