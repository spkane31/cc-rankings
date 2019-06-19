from bs4 import BeautifulSoup
import requests
import csv
import os
import numpy as np
import json


TFRRS_TEST_RACE = 'https://www.tfrrs.org/results/xc/14671/Roy_Griak_Invitational'

def getLinks(URL):
  headers = {'User-Agent': 'Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.76 Safari/537.36'} # This is chrome, you can set whatever browser you like
  
  s = requests.session()
  
  r = s.get(URL, headers=headers)
  soup = BeautifulSoup(r.content, 'html.parser')
  # <div class="col-lg-12">
  races = soup.find_all('div', class_='row')#, class_='col-lg-12')
  
  for race in races:
    details = {}
    name = race.find('h3', class_='font-weight-500')
    try:
      info = name.get_text()
      info = info.split()
      if len(info) > 0:
        for i in info:
          if str.upper(i) in ['MEN', 'M', 'MENS', '(M)']:
            details['gender'] = 'MENS'
          if str.upper(i) in ['WOMEN', 'W', 'WOMENS', '(W)']:
            details['gender'] = 'WOMENS'
          if str.upper(i) in ['INDIVIDUAL']:
            details['valid'] = True
          elif str.upper(i) in ['TEAM']:
            details['valid'] = False
          if str.upper(i) in ['8K', '6K', '5K', '5 MILE']:
            details['distance'] = str.upper(i)

      results = race.find('tbody', class_='color-xc')
      print(len(results))

      getResults(results)

      if details['valid']: print(details)

    except:
      pass

def getResults(results):
  results = results.find_all('tr')
  results = [r.get_text() for r in results]
  results = [str(r) for r in results]
  results = [r.split('\n') for r in results]
  results = [[r[3], r[6], r[9], r[15]] for r in results]
  print(results)
  
  quit()

getLinks(TFRRS_TEST_RACE)