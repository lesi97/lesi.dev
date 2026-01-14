create index if not exists idx_destiny_user_characters_membership_id
on public.destiny_user_characters (membership_id);

create index if not exists idx_destiny_weapon_kill_counts_membership_weapon_hash
on public.destiny_weapon_kill_counts (membership_id, weapon_hash);

create index if not exists idx_destiny_weapon_kill_counts_membership_weapon_id
on public.destiny_weapon_kill_counts (membership_id, weapon_id);

create index if not exists idx_destiny_weapon_kill_counts_weapon_hash
on public.destiny_weapon_kill_counts (weapon_hash);

create index if not exists idx_destiny_weapons_bucket_type_hash
on public.destiny_weapons (bucket_type_hash);

drop table user_profiles;
drop view terror_goal_messages;
drop table terror_goal;
drop trigger on_auth_user_created on auth.users;
drop function if exists "public"."handle_new_user";
drop function if exists "public"."update_updated_at_column";