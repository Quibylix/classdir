ALTER TABLE presentations ADD COLUMN content TEXT NOT NULL DEFAULT '';

UPDATE presentations p
SET content = COALESCE(
  (SELECT string_agg(s.content, E'\n---\n' ORDER BY array_position(p.slide_order, s.id))
   FROM slides s
   WHERE s.presentation_id = p.id),
  ''
);

ALTER TABLE presentations DROP COLUMN slide_order;

DROP TABLE slides;
