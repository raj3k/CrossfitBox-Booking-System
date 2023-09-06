import { Login, Register, Layout } from '@/views/account';

export default {
  path: '/account',
  component: Layout,
  children: [
    { path: '', redirect: '/account/login' },
    { path: 'login', component: Login },
    { path: 'register', component: Register }
  ]
}