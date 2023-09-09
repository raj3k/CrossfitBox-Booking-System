interface User {
  id: string,
  email: string,
  first_name: string,
  last_name: string,
  is_active: boolean,
  is_staff: boolean,
  is_superuser: boolean,
  thumbnail: string,
  created_at: number,
  profile: Profile
}

interface Profile {
  id: string,
  user_id: string,
  phone_number: number,
  birth_date: number,
}