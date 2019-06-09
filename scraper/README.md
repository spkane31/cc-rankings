# Scraper - Go
This web scraper will pull all results from  a TFRRS link and convert the page
into a mens and womens csv with the results and a JSON file with details about
the race. The format of the mens/womens CSV is:
```
<last>, <first>, <school>, <year>, <time>
```

## Issues
Some of the results pages will list the team results for multiple division, ie (https://www.tfrrs.org/results/xc/15030/NCCAA_National_Championships), which throws off how the data is scraped. For now I'm ignoring these results, but will come back to it, there's a lot of races that have multiple division scoring in them. 