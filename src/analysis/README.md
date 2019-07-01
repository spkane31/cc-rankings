# Data Analysis


## Conversions Table
Each point increase in rating, equates to approximately 1 second/mile of race distance. A one point difference in the men's conversions, which are standardized
for 8K races, is roughly 5 seconds per mile, and for the women standardized to the 5K, is 3 seconds per mile.

The formula boils down to:
( base - actual - correction ) / (distance in miles)

For men: (1900 - actual - correction ) / (8 / 1.609)
For women: (1350 - actual - correction ) / (6 / 1.609)

The goal is to have the ratings scale from 0 to 100. With zero being a middle of the pack club runner, and 100 being national title contender. The baseline is a very quick time for a neutral cross country course (ie. rolling gentle hills, long straightaways, minimum narrow points, etc.). This is roughly 23:23 for men's 8k (29:38 for 10k) and 16:17 for women's 5k (19:45 for 6k). I want these to be roughly equivalent to keep the rating system simple as both of these times equate to about 90.7 percent age graded time for a 22 year old runner. 

| Rating | Men's Time | Women's Time |
| ------ |:----------:|:------------:|
| 100    | 23:23      | 16:17        |
| 90     | 24:13      | 16:54        |
| 80     | 25:02      | 17:32        |
| 70     | 25:52      | 18:09        |
| 60     | 26:42      | 18:46        |
| 50     | 27:31      | 19:24        |
| 40     | 28:21      | 20:01        |
| 30     | 29:11      | 20:38        |
| 20     | 30:01      | 21:15        |
| 10     | 30:50      | 21:53        |
|  0     | 31:40      | 22:30        |
