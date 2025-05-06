export const ResponseMap = {
    Success: {
        ret_code: 0,
        message: 'success'
    },
    UserExist: {
        ret_code: 10000,
        message: 'username already exist'
    },
    SystemError: {
        ret_code: 10001,
        message: 'system error'
    },
    UserError: {
        ret_code: 10002,
        message: 'username or password error'
    },
    ParamsError: {
        ret_code: 10003,
        message: 'params error'
    },
    EmailError: {
        ret_code: 10004,
        message: 'email already exist'
    },
    LoginError: {
        ret_code: 10005,
        message: 'user not login'
    },
    AuthCodeError: {
        ret_code: 10006,
        message: 'Incorrect verification code.'
    },
    TokenExpired: {
        ret_code: 10009,
        message: 'The login credentials have expired.'
    },
    SystemEmptyError: {
        ret_code: 10010,
        message: 'No data found.'
    },
    RoleExisted: {
        ret_code: 10011,
        message: 'The role or name cannot be repeated in the same namespace.'
    }
}

export const [HttpCodeSuccess, HttpCodeNotFound] = [200, 400]

export const ResourceCategory = ["API", "Menu", "Action", "Other"];