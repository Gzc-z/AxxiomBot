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

| impl           | status          |
|--------------- | --------------- |
| Item1.2        | Pending         |
| Item1.3        | Pending         |
| clima        | clima           |
| minigames      | coming soon     |


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

# nome

eu sinceramente não sei que nome dar a esse bot

# checklist

* [x] escalar o bot para diferentes tags com diferentes propositos
* [x] atualizar README e documentações
* [ ] refatorar o código
* [ ] colocar mais coisas na checklist
* [ ] escolher um bom nome
* [ ] fazer o bot multilinguagem
