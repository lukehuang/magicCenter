"LinkTest"

import MagicSession
import LoginTest

def join_str(catalog):
    'JoinStr'
    ret = ''
    for v in catalog:
        ret = '%s,%d'%(ret, v)
    return ret

class LinkTest(MagicSession.MagicSession):
    'LinkTest'
    def __init__(self, auth_token):
        MagicSession.MagicSession.__init__(self)
        self.authority_token = auth_token

    def create(self, name, url, logo, catalogs):
        'create'
        params = {'link-name': name, 'link-url': url, 'link-logo': logo, 'link-catalog': catalogs}
        val = self.post('http://localhost:8888/content/link/?token=%s'%self.authority_token, params)
        if val and val['ErrCode'] == 0:
            print 'create link success'
            return val['Link']

        print 'create link failed'
        return None

    def destroy(self, link_id):
        'destroy'
        val = self.delete('http://localhost:8888/content/link/%s/?token=%s'%(link_id, self.authority_token))
        if val and val['ErrCode'] == 0:
            print 'destroy link success'
            return True

        print 'destroy link failed'
        return False

    def update(self, link):
        'update'
        catalogs = join_str(link['Catalog'])
        params = {'link-name': link['Name'], 'link-url': link['URL'], 'link-logo': link['Logo'], 'link-catalog': catalogs}
        val = self.put('http://localhost:8888/content/link/%s/?token=%s'%(link['ID'], self.authority_token), params)
        if val and val['ErrCode'] == 0:
            print 'update link success'
            return val['Link']

        print 'update link failed'
        return None

    def query(self, link_id):
        'query'
        val = self.get('http://localhost:8888/content/link/%d/?token=%s'%(link_id, self.authority_token))
        if val and val['ErrCode'] == 0:
            print 'query link success'
            return val['Link']

        print 'query link failed'
        return None

    def query_all(self):
        'query_all'
        val = self.get('http://localhost:8888/content/link/?token=%s'%self.authority_token)
        if val and val['ErrCode'] == 0:
            print 'query_all link success'
            return val['Link']

        print 'query_all link failed'
        return None

if __name__ == '__main__':
    LOGIN = LoginTest.LoginTest()
    if not LOGIN.login('rangh@126.com', '123'):
        print 'login failed'
    else:
        APP = LinkTest(LOGIN.authority_token)
        LINK = APP.create('testLink', 'test link url', 'test link logo', '8,9')
        if LINK:
            LINK_ID = LINK['ID']
            LINK['URL'] = 'aaaaaa, bb dsfsdf  erewre aa'
            LINK['Logo'] = 'test link logo'
            LINK['Catalog'] = [8,9,10]
            LINK = APP.update(LINK)
            if not LINK:
                print 'update link failed'
            elif len(LINK['Catalog']) != 3:
                print 'update link failed, link len invalid'
            else:
                pass
            LINK = APP.query(LINK_ID)
            if not LINK:
                print 'query link failed'
            elif cmp(LINK['URL'],'aaaaaa, bb dsfsdf  erewre aa') != 0:
                print 'update link failed, content invalid'

            if len(APP.query_all()) <= 0:
                print 'query_all link failed'
            
            APP.destroy(LINK_ID)
        else:
            print 'create link failed'

        LOGIN.logout(LOGIN.authority_token)