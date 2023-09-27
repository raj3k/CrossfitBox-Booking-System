import axios from "axios";
import { UserData } from "../types/user";

const API_SERVER = "http://localhost:4000"
const instance = axios.create({
  withCredentials: true,
  baseURL: API_SERVER,
});

export function signup(email: string, firstName: string, lastName: string, password: string) {
  return instance.post<UserData>("/api/v1/users/register/", {
    email: email,
    first_name: firstName,
    last_name: lastName,
    password: password,
  });
}

export function activateUser(userId: string, token: string) {
  return instance.put(`/api/v1/users/activate/${userId}`, {
    token: token.split(' ').join('')
  });
}

export function signIn(email: string, password: string) {
  return instance.post<UserData>("/api/v1/users/login/", {
    email: email,
    password: password,
  }, {
    withCredentials: true
  });
}

export function getCurrentUser() {
  return instance.get<UserData>("/api/v1/users/current-user", {
    withCredentials: true
  });
}

export function signOut() {
  return instance.post("/api/v1/users/logout");
}