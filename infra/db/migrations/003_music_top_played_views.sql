create or replace view music.top_played_songs with (security_invoker = on) as
select
  a.name as artist,
  al.name as album,
  t.name as track,
  count(*)::bigint as listened_count,
  min(s.scrobbled_at) as first_listened,
  max(s.scrobbled_at) as last_listened,
  t.url as track_url,
  coalesce(t.image_url, al.image_url) as track_art
from music.scrobbles s
join music.artists a on a.id = s.artist_id
join music.tracks t on t.id = s.track_id
left join music.albums al on al.id = s.album_id
group by
  a.id,
  a.name,
  al.id,
  al.name,
  al.image_url,
  t.id,
  t.name,
  t.url,
  t.image_url
order by
  listened_count desc,
  last_listened desc,
  artist,
  track;

create or replace view music.top_played_albums with (security_invoker = on) as
select
  a.name as artist,
  al.name as album,
  count(*)::bigint as listened_count,
  min(s.scrobbled_at) as first_listened,
  max(s.scrobbled_at) as last_listened,
  al.url as album_url,
  al.image_url as album_art
from music.scrobbles s
join music.artists a on a.id = s.artist_id
join music.albums al on al.id = s.album_id
group by
  a.id,
  a.name,
  al.id,
  al.name,
  al.url,
  al.image_url
order by
  listened_count desc,
  last_listened desc,
  artist,
  album;

create or replace view music.top_played_artists with (security_invoker = on) as
select
  a.name as artist,
  count(*)::bigint as listened_count,
  min(s.scrobbled_at) as first_listened,
  max(s.scrobbled_at) as last_listened,
  a.url as artist_url,
  a.image_url as artist_art
from music.scrobbles s
join music.artists a on a.id = s.artist_id
group by
  a.id,
  a.name,
  a.url,
  a.image_url
order by
  listened_count desc,
  last_listened desc,
  artist;

comment on view music.top_played_songs is 'Aggregated top played songs for Grafana music dashboards.';
comment on view music.top_played_albums is 'Aggregated top played albums for Grafana music dashboards.';
comment on view music.top_played_artists is 'Aggregated top played artists for Grafana music dashboards.';
