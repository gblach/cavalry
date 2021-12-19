import os, unittest
import requests

connect_to = os.getenv('CONNECT_TO')

class Requests(unittest.TestCase):
    def test_root_status(self):
        r = requests.get(connect_to + '/')
        self.assertEqual(r.status_code, 404)

    def test_quote_status(self):
        r = requests.get(connect_to + '/quote')
        self.assertEqual(r.status_code, 200)

    def test_quote_json(self):
        r = requests.get(connect_to + '/quote')
        json = r.json()
        self.assertEqual(json['quote'], "I've always been more interested in the future than in the past.")
        self.assertEqual(json['author'], "Grace Hopper")
