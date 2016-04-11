# Matlistan

Slumpa måltider för valt antal dagar och få en sammanställd inköpslista.

Projektet är skrivet i Javascript/jQuery och använder sig av en statisk
"matdatabas" i JSON-format, vilken man får uppdatera manuellt. Hur denna fil
genereras är upp till användaren själv. I nuläget finns det ett medföljande
python-skript som omvandlar "Matlistan.txt"-formatet till för projektet rätt
JSON-syntax.

Om du har Linux/Mac OS X kan du köra skriptet i terminal:

> ./create-json-data.py < Matlista.txt > data.json

Projektet är under utveckling. Se mer under [TODO.md](./TODO.md).
