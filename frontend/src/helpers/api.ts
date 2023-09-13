import axios from "axios";
// TODO: resolve http://localhost:4000/

export function signup(email: string, firstName: string, lastName: string, password: string) {
  return axios.post<User>("http://localhost:4000/v1/users/register/", {
    email: email,
    first_name: firstName,
    last_name: lastName,
    password: password,
  });
}

export function activateUser(userId: string, token: string) {
  return axios.put(`http://localhost:4000/v1/users/activate/${userId}`, {
    token: token.split(' ').join('')
  });
}