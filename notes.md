dans world:
ajouter une map[int]map[int]WorldChunk (matrice creuse)

avoir une méthode pour savoir quand est-ce qu'un chunk doit être rendu ou non
en fonction de la taille du viewport et de son centre

dans la boucle de jeu, quand le vaisseau approche est dans un chunk, tous les chunks adjacents sont générés si ils ne l'étaient pas déjà

un chunk est constitué de 32*32 cellules.
chaque cellule fait 50px de côté.
