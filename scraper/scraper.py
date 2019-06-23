# import pprint
from bs4 import BeautifulSoup
import requests
import csv
import os
import numpy as np
import json, time, logging
import time

logging.basicConfig(format='%(asctime)s - %(message)s', level=logging.INFO)

NIRCA_RACES = 'https://clubrunning.org/races/?season=F-18'
TFRRS_RACES = 'https://www.tfrrs.org/results_search.html'
NIRCA_TEST_RACE = 'https://clubrunning.org/races/race_results.php?race=677'
TFRRS_TEST_RACE = 'https://www.tfrrs.org/results/xc/15028/NCAA_Division_III_Cross_Country_Championships'

def removeSpecialCharacters(s):
  s = s.replace(' ', '')
  s = s.replace('/', '')
  s = s.replace(',', '')
  s = s.replace(')', '')
  s = s.replace('(', '')
  s = s.replace('*', '')
  s = s.replace("'", '')
  s = s.replace('"', '')
  s = s.replace('|', '')
  s = s.replace('&', '')
  return s

def removeSpecialCharactersNotSpaces(s):
  s = s.replace('/', '')
  s = s.replace(',', '')
  s = s.replace(')', '')
  s = s.replace('(', '')
  s = s.replace('*', '')
  s = s.replace("'", '')
  s = s.replace('"', '')
  s = s.replace('|', '')
  s = s.replace('&', '')
  return s


def getNIRCALinks(URL):
    s = requests.session()
    BASE = 'https://clubrunning.org/races/'
    r = s.get(URL)
    soup = BeautifulSoup(r.content, 'html.parser')

    races = soup.find_all('tr', class_='racerow')
    allRaces = []
    for r in races:
        temp = []
        details = r.find_all('td', class_='column1a row')
        details = [d.get_text() for d in details]
        
        date = details[0].split('2018')
        date = date[0] + '2018'
        
        location = details[0].split('Hosted')
        location = location[0].split('2018')[1]

        link = r.find_all('a', href=True)
        race = link[0].get_text()
                
        link = BASE + link[0]['href']
        link = link.replace('info', 'results')

        temp = [link, race, date, location]
        allRaces.append(temp)

    # allRaces = [link, date, location]
    return allRaces

def getTFRRSLinks(month=None, year=None, meet_name=None, state=None):
  headers = {'User-Agent': 'Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.76 Safari/537.36'} # This is chrome, you can set whatever browser you like
  data = {'meet_name': '', 'sport': 'xc', 'state': '', 'month': '', 'year': ''}
  if month != None:
    data['month'] = str(month)
  if year != None:
    data['year'] = str(year)
  if meet_name != None:
    data['meet_name'] = meet_name
  if state != None:
    data['state'] = state

  s = requests.session()
  
  r = s.get(TFRRS_RACES, headers=headers, params=data)
  soup = BeautifulSoup(r.content, 'html.parser')
  races = soup.find_all('tbody', id='results_page1')#, href=True)
  races = races[0].find_all('a')
  links = [r['href'][:-1] for r in races]     # need the [:-2] to get rid of the '\n' characters
  
  not_except = True
  count = 2
  while not_except:
    try:
      string = 'results_page' + str(count)
      races = soup.find_all('tbody', id=string)#, href=True)
      races = races[0].find_all('a')
      l = [r['href'][:-1] for r in races]     # need the [:-2] to get rid of the '\n' characters
      if len(l) == 0: not_except = False
      links += l
      count += 1

    except:
      not_except = False
  print(f"Found {len(links)} Links!")
  return links

def getNIRCAResults(URL):
    s = requests.session()
    r = s.get(URL)
    soup = BeautifulSoup(r.content, 'html.parser')
    races = soup.find_all('tbody')

    womens = []
    mens = []
    raceBool = False
    count = 0
    
    raceName= soup.find_all('title')
    raceName = soup.find_all('span', class_='style1')
    raceName = [r.get_text() for r in raceName]
    raceName = [str(r) for r in raceName]
    raceName = raceName[0].replace(' ', '')

    titles = soup.find_all('span', class_='style2')
    titles = [title.get_text() for title in titles]
    titles = [str(title) for title in titles]

    for r in races:
        performances = r.find_all('tr')
        performances = [p.get_text().replace(',','') for p in performances]
        racers = [str(p) for p in performances]
        racers = [r.split('\n') for r in racers]
        try:
            results = [[x[3], x[4], x[5], x[6], x[8]] for x in racers]
            raceBool = True
            count += 1
        except:
            raceBool = False
        if count % 2 == 1 and raceBool:
            mens.append(results)
        elif count % 2 == 0 and raceBool:
            womens.append(results)
    
    mens = np.asarray(mens)
    womens = np.asarray(womens)
    return mens, womens #,  raceName

def getTFRRSResults(URL):

  headers = {'User-Agent': 'Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.76 Safari/537.36'} # This is chrome, you can set whatever browser you like
  
  s = requests.session()
  
  r = s.get(URL, headers=headers)
  soup = BeautifulSoup(r.content, 'html.parser')
  races = soup.find_all('div', class_='row')
  
  date_loc = soup.find_all('div', class_='panel-heading-normal-text inline-block')
  if not len(date_loc):
    print('Date did not work')
    return

  date = date_loc[0].get_text().replace("\n", " ")
  date = date.strip()
  loc = date_loc[1].get_text().replace("\n", " ")
  loc = loc.strip()

  race_name = soup.find_all('a', class_='white-underline-hover')
  race_name = race_name[0].get_text().replace("\n", " ")
  race_name = race_name.replace("\t", " ")
  race_name = race_name.strip()

  print(race_name)

  m = []

  for race in races:
    details = {}
    name = race.find('h3', class_='font-weight-500')
    if name != None:
    # try:
      details['date'] = date
      loc = removeSpecialCharactersNotSpaces(loc)
      details['course'] = loc
      details['name'] = race_name
      info = name.get_text()
      info = info.split()
      info = [str.upper(u) for u in info]
      info = [i.replace('(', '') for i in info]
      if len(info) > 0:
        new_race_name = ""
        finished = False
        for i in info:
          if str.upper(i) in ['MEN', 'M', 'MENS', '(M)', "MEN'S"]:
            details['gender'] = 'MENS'
          if str.upper(i) in ['WOMEN', 'W', 'WOMENS', '(W)', "WOMEN'S"]:
            details['gender'] = 'WOMENS'
          if str.upper(i) in ['INDIVIDUAL']:
            details['valid'] = True
          elif str.upper(i) in ['TEAM']:
            details['valid'] = False
          if str.upper(i) in ['8K', '6K', '4K', '3K', '5K', '5', '10K', '10000', '6000', '4000', '8000', '5000', '3200' '3.73', '3', '4.97', '4.96' '3.1', '4', '8.4', '2.95', '7900', '7K', '8369']:
            if str.upper(i) == '10000':
              details['distance'] = '10K'
            elif str.upper(i) == '6000':
              details['distance'] = '6K'
            elif str.upper(i) == '4000':
              details['distance'] = '4K'
            elif str.upper(i) == '4K':
              details['distance'] = '4K'
            elif str.upper(i) == '8000':
              details['distance'] = '8K'
            elif str.upper(i) == '5000':
              details['distance'] = '5K'
            elif str.upper(i) == '5':
              if 'MILE' in info:
                details['distance'] = '5 MILE'
              elif 'K' in info:
                details['distance'] = '5K'
            elif str.upper(i) == '3K':
              details['distance'] = '3K'
            elif str.upper(i) == '4.97':
              details['distance'] = '8K'
            elif str.upper(i) == '4.96':
              details['distance'] = '8K'
            elif str.upper(i) == '2.95':
              details['distance'] = '2.95'
            elif str.upper(i) == '3.73':
              details['distance'] = '6K'
            elif str.upper(i) == '3.1':
              details['distance'] = '5K'
            elif str.upper(i) == '3200':
              details['distance'] = '3.2K'
            elif str.upper(i) == '8.4':
              details['distance'] = '8.4K'
            elif str.upper(i) == '8369':
              details['distance'] = '8.369K'
            elif str.upper(i) == '3':
              if 'MILE' in info:
                details['distance'] = '3 MILE'
            elif str.upper(i) == '4':
              if 'MILE' in info:
                details['distance'] = '4 MILE'
            elif str.upper(i) == '7900':
              details['distance'] = '7.9K'
            elif str.upper(i) == '7K':
              details['distance'] = '7K'
            else:
              details['distance'] = str.upper(i)

          if str.upper(i) in ["RESULTS", "RESULT", "INDIVIDUAL"]:
            finished = True
          elif not finished:
            new_race_name += i + " "
      results = race.find('tbody', class_='color-xc')

      if details['valid']: 
        r = getResults(results)
        details['race_name'] = new_race_name
        details['results'] = r
        m.append(details)
        try:
          a = details['distance']
        except:
          logging.warning(f"No Distance found at: {race_name}, {date}")
          # print("\t\tThis didn't work!!! No Distance")
          return
        try:
          a = details['gender']
        except:
          logging.warning(f"No gender found at: {race_name}, {date}")
          # print("\t\tThis didn't work!!! No gender")
          return
  return write_results(m)
  
  # quit()

def getResults(results):
  results = results.find_all('tr')
  results = [r.get_text() for r in results]
  results = [str(r) for r in results]
  results = [r.split('\n') for r in results]
  # results = [[r[3], r[6], r[9], r[15]] for r in results]
  ret = []
  for result in results:
    names = result[3].split(",")
    year = result[6]
    school = result[9]
    time = result[15]
    names[0] = names[0].strip()
    names[1] = names[1].strip()

    names[0] = names[0].replace("'", "")
    names[1] = names[1].replace("'", "")
    
    names[0] = names[0].replace('"', "")
    names[1] = names[1].replace('"', "")

    names[0] = names[0].replace(",", "")
    names[1] = names[1].replace(",", "")
    names[0] = names[0].replace(" ", "")
    names[1] = names[1].replace(" ", "")
    ret.append([names[0], names[1], year, school, time])
  return ret
  
def write_results(m):
  if not len(m): return
  count = 0
  directory = 'RaceResults/'
  json_data = {}

  json_data['date'] = m[0]['date']
  json_data['course'] = m[0]['course']
  json_data['name'] = m[0]['name']

  try:
    os.mkdir(directory)
  except:
    pass
    
  race = m[0]['name']
  race = removeSpecialCharacters(race)
  race = race.upper()
  
  year = json_data['date'].split()[-1]
  year = removeSpecialCharacters(year)

  directory = os.path.join(directory, race)

  try:
      os.mkdir(directory)
  except:
    pass


  directory = os.path.join(directory, year)

  try:
      os.mkdir(directory)
  except:
    pass


  index = 0
  for race in m:
    index += 1
    file_name = 'file' + str(index)
    json_data[file_name] = {}
    json_data[file_name]['race_name'] = race['race_name']
    json_data[file_name]['gender'] = race['gender']
    json_data[file_name]['distance'] = race['distance']

    results_file_name = race['race_name'].replace(' ', '')
    results_file_name = removeSpecialCharacters(results_file_name)
    results_file_name = results_file_name.upper()

    file = os.path.join(directory, results_file_name+'.csv')
    try:
      new_file = (results_file_name + '.csv')
      new_file = removeSpecialCharacters(new_file)
      json_data[file_name]['file'] = new_file
      count += len(race['results'])
      np.savetxt(os.path.join(directory, new_file), race['results'], delimiter=", ", fmt="%s")
    except:
      pass

  with open(directory+'/raceSummary.json', 'w') as f:
    json.dump(json_data, f)

  return count

start = time.time()
allRaces = getTFRRSLinks(meet_name='NCAA Division I')
count = 0
for race in allRaces:
  
  a = getTFRRSResults('http:'+race)
  if a != None:
    count += a
    # print(f"Count = {count}\tTime = {time.time() - start}\tAverage = {count / (time.time() - start)}")

print(f"Found {count} results in {time.time() - start} seconds! Average {count / (time.time() - start)} per second!")