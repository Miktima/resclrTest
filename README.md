# resclrTest
Golang script for test links throughth CDN to image rescaler. The goal is to test image rescaler by loading 10 factor requests.  

Incoming parameters:
- URL page
- CDN address

The parameters may be loaded direct in CLI or throughth conf.json file. 

URLs to webp or jpg images are extracted from the URL page. The CDN address is used for this extraction. This CDN address can use regular expression for multiple CDN addresses,
for example, cdn(1|2).example.com. Then the extracted images are requested asynchronously with factor 10 (qnty images * 10). If response status is not ok (200) then error message occurs, 
otherwise dot (.) is printed.
