import axios from "axios";

export function signup(email: string, firstName: string, lastName: string, password: string) {
  return axios.post<User>("http://localhost:4000/v1/users/register/", {
    email: email,
    first_name: firstName,
    last_name: lastName,
    password: password,
  });
}