# AxxiomBot

Um bot administrativo para o Discord, projetado para facilitar processos, tarefas e problemas comuns.  
futuramente pretendo ligar esse bot a outros contextos, se estendendo com outros projetos.  

funcionalidades atuais:
* criar um grupo de pontuações (tags)
* listagem de perfil de usuários
* pegar fatos de catfact (https://catfact.ninja)
* pegar tempo de resposta de um servidor
* calculadora básica

implementações futuras:

- minigames
- meme aleatório
- jogo da vida (de alguma forma)
- 1 musica aleatória por dia
- jogo de palavras
- rank de algo com enquete

# configuração

[discord applications](https://discord.com/developers/applications)

* primeiramente é preciso colocar informações 
* primeiramente acesse o discord application;
* crie um novo aplicativo
* acesse a aba bot e copie o token resetando ele:  
* cole o token e preencha outras informações em `.env_example` na base do projeto
* renomeie para `.env`

em seguida use `make run` na base do projeto para iniciar o bot

criar novo aplicativo:  
[![criar nova aplicação](./docs/New_app.png)](https://discord.com/developers/applications#:~:text=a)
![redefinir e copiar token](./docs/token.png)

# nome e contexto

# checklist

* [x] escalar o bot para diferentes tags com diferentes propositos
* [x] atualizar README e documentações
* [x] colocar mais coisas na checklist
* [x] escolher um bom nome
* [ ] colocar mais coisas nas implementações futuras
* [ ] refatorar o código gambiarrento
* [ ] fazer o bot para outros idiomas
