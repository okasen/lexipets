-- Create a keyspace
CREATE KEYSPACE IF NOT EXISTS lexipets WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : '1' };

-- Create types
CREATE TYPE IF NOT EXISTS lexipets.feature (
    part text,
    dominant text,
    recessive text,
    mixed text
                               );

CREATE TYPE IF NOT EXISTS lexipets.gene (
    dominant boolean,
    recessive boolean,
    features list <frozen <feature>>
                               );

-- Create a table
CREATE TABLE IF NOT EXISTS lexipets.species (
                                                id timeuuid PRIMARY KEY,
                                                name text,
                                                features list <frozen <lexipets.feature>>
);

CREATE TABLE IF NOT EXISTS lexipets.pets (
                                             owner_id timeuuid,
                                             id timeuuid,
                                             name text,
                                             species_name text,
                                             species_features list <frozen <lexipets.feature>>,
                                             genes list <frozen <lexipets.gene>>,
                                             img text,
                                             PRIMARY KEY(owner_id, id)
);

CREATE TABLE IF NOT EXISTS lexipets.users (
    id timeuuid PRIMARY KEY,
                                          username text,
                                          password text,
    email text,
                                          created_at timestamp,
                                          updated_at timestamp,
                                          deleted_at timestamp
);

-- Insert some data
INSERT INTO lexipets.species
(id, name, features)
VALUES (now(), 'Waggler', [{
        part: 'ears',
        dominant: 'pointy',
        recessive: 'floppy',
        mixed: 'floppy'
        },
    {
        part: 'coatColor',
        dominant: 'brown',
        recessive: 'yellow',
        mixed: 'brown'
        },
    {
        part: 'tail1',
        dominant: 'curled',
        recessive: 'straight',
        mixed: 'curved'
        },
    {
        part: 'tail2',
        dominant: 'tight',
        recessive: 'loose',
        mixed: 'relaxed'
        }
    ]);
INSERT INTO lexipets.species
(id, name, features)
VALUES (now(), 'Gooper', [{
    part: 'jelly',
    dominant: 'firm',
    recessive: 'fluid',
    mixed: 'viscous'
    },
    {
        part: 'pattern',
        dominant: 'splotches',
        recessive: 'plain',
        mixed: 'speckles'
        },
    {
        part: 'horn',
        dominant: 'present',
        recessive: 'missing',
        mixed: 'missing'
        },
    {
        part: 'hornstyle',
        dominant: 'corkscrew',
        recessive: 'straight',
        mixed: 'ridged'
        }
    ]);