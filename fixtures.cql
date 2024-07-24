-- Create a keyspace
CREATE KEYSPACE IF NOT EXISTS lexipets WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : '1' };

-- Create a table
CREATE TABLE IF NOT EXISTS lexipets.species (
                                                id text PRIMARY KEY,
                                                name text,
                                                features list <frozen <map<text, text>>>
);
CREATE TABLE IF NOT EXISTS lexipets.pets (
                                             id text PRIMARY KEY,
                                             name text,
                                             species list <frozen <map<text, text>>>,
                                             genes list <frozen <map<text, text>>>,
                                             img text
);

-- Insert some data
INSERT INTO lexipets.species
(id, name, features)
VALUES ('1', 'Waggler', [{
        'part': 'ears',
        'dominant': 'pointy',
        'recessive': 'floppy',
        'mixed': 'floppy'
        },
    {
        'part': 'coatColor',
        'dominant': 'brown',
        'recessive': 'yellow',
        'mixed': 'brown'
        },
    {
        'part': 'tail1',
        'dominant': 'curled',
        'recessive': 'straight',
        'mixed': 'curved'
        },
    {
        'part': 'tail2',
        'dominant': 'tight',
        'recessive': 'loose',
        'mixed': 'relaxed'
        }
    ]);
INSERT INTO lexipets.species
(id, name, features)
VALUES ('2', 'Gooper', [{
    'part': 'jelly',
    'dominant': 'firm',
    'recessive': 'fluid',
    'mixed': 'viscous'
    },
    {
        'part': 'pattern',
        'dominant': 'splotches',
        'recessive': 'plain',
        'mixed': 'speckles'
        },
    {
        'part': 'horn',
        'dominant': 'present',
        'recessive': 'missing',
        'mixed': 'missing'
        },
    {
        'part': 'hornstyle',
        'dominant': 'corkscrew',
        'recessive': 'straight',
        'mixed': 'ridged'
        }
    ]);