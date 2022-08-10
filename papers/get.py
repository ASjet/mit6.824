import requests
import bs4
import re
import html5lib

url = "https://pdos.csail.mit.edu/6.824/schedule.html"
baseurl = "https://pdos.csail.mit.edu/6.824/"
ptn = re.compile(".*\.pdf$")
resp = requests.get(url)
content = resp.content.decode("utf-8")

html = bs4.BeautifulSoup(content, "html5lib")
tags = html.body.findAll("a")

links = [baseurl + tag.attrs["href"]+'\n' for tag in tags if ptn.match(tag.attrs["href"])]
for link in links:
    print(link)

