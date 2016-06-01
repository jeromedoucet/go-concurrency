package client_test

// first step : lauch the client (and start the mock server !)
// second step : register the server (consul ?, zk ? inner ?) and check the registration + have a health check ?
// third step : receive the request
// fourth step : reply
// finally : check in redis that the order has been deleted and that the score has been incremented

/*
 apres enregistrement d'un participant, on commence a lui envoyer des commandes en mode synchrone. pas de mux dyn je crois :(
 si un echec, on ballance des health check pendant 10 minutes jusqu'a ce qu'il y en ait 1 qui passe.
 Si c'est le cas, on recommence, sinon on supprime.

 Si un enregistrement se produit plusieurs fois de suite (meme couple ip / nom) -> idempotence
 Si tentative d'enregistrement d'une ip existant deja avec un autre nom => erreur
 Si tentative d'enregistrenebt d'une nouvelle ip avec nom deja existant => erreur
*/
