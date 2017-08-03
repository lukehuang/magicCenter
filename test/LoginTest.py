'LoginTest'

import MagicSession

class MagicCenter(MagicSession.MagicSession):
    'MagicCenter'
    def __init__(self):
        MagicSession.MagicSession.__init__(self)
        self.authority_token = ''
        self.current_user = None

    def login(self, account, password):
        'login'
        params = {'user-account': account, 'user-password': password}
        val = self.post('http://localhost:8888/cas/user/', params)
        if val and val['ErrCode'] == 0:
            self.authority_token = val['AuthToken']
            self.current_user = val['User']
            print 'login success'
            print self.current_user
            print self.authority_token
            return True
        else:
            print 'login failed'
            return False

    def logout(self, authToken):
        'logout'
        val = self.delete('http://localhost:8888/cas/user/?authToken=%s'%authToken)
        if val and val['ErrCode'] == 0:
            print 'logout success'
        else:
            print 'logout failed'

if __name__ == '__main__':
    APP = MagicCenter()
    APP.login('rangh@126.com', '123')
    APP.logout('123')
    APP.logout(APP.authority_token)