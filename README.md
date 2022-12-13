# SmartBot
É um robô que monitora a cotação de ativos (ações, etfs e smartcoins) e alerta quando 
atingir um determinado limiar de ganhos ou perdas.  Além disso, informa se você deve 
vender ou comprar ativos, monitorando e analisando o RSI*.

> RSI é um indicador técnico utilizado na análise do mercado. Ele rastreia a tendência de 
> preço de um ativo. Para saber mais: [RSI](https://www.suno.com.br/artigos/rsi/)

## Atenção
Esse programa não garante lucros nem impede perdas financeiras. Ele foi criado para fins
didáticos. Use por sua conta e risco.

## Funcionamento
Dada uma carteira de ativos é feita, num intervalo de tempo configurável, a leitura das 
informações necessárias para notificar o usuário quando atingir um valor pré-determinado de 
ganho ou perda. Da mesma forma, quando o RSI do ativo for menor que 30 pontos, o sistema 
indica a compra, e se for maior que 70 pontos, indica a venda. A informação é enviada através
de uma mensagem via Telegram e/ou email.

![smartbot-cli](https://user-images.githubusercontent.com/12521070/190633097-2381e834-1fff-46c2-97fa-53ff3fb54f22.png)

![Imagem real das mensagem no Telegram](https://user-images.githubusercontent.com/12521070/190629524-eb12384a-a7c5-46b5-a2ea-a677533d1ab7.jpg)

###### Imagem real das mensagens enviadas ao Telegram.


## Configuração
Existem 2 arquivos de configuração: smartbot.cfg e carteira.cfg (no repositório oo modelos 
estão com a extensão .sample, basta renomea-los).

### smartbot.cfg
Cuida da configuração do programa em si.

    {  
       "SleepMinutes": 5,  
       "telegramID": 47283687,  
       "telegramToken": "1237124367:AHUs7HPOiHTRdthRRE7VC4CES4d5D44DA01",  
       "emailLogin": "seuemail@gmail.com",  
       "emailPassword": "suasenha123456",  
       "emailTo": "outroemail@gmail.com"  
    }
**SleepMinutes**: informa o tempo em minutos que o programa deve aguardar para realizar as
buscas por novas cotações. Sugiro não colocar um valor muito baixo para evitar o bloqueio 
do IP.

**telegramID**: é o seu ID no Telegram. Para descobrir qual é, veja aqui:
[How to get an id to use on Telegram Messenger](https://github.com/GabrielRF/telegram-id). 
Se você não quiser utilizar o Telegram, basta colocar 0 (zero) neste campo.

**telegramToken**: é o token do bot do Telegram que irá receber as mensagens de alerta.
Existem diversos tutoriais ensinando a fazer isso, como 
[10 Passos para se criar um Bot no Telegram](https://medium.com/tht-things-hackers-team/10-passos-para-se-criar-um-bot-no-telegram-3c1848e404c4).

**emailLogin**: login da conta que fará o envio do email.

**emailPassaword**: senha da conta que fará o envio do email.

**emailTo**: endereço de email que receberá as mensagens do smartbot.

### carteira.cfg
Para elaboração da carteira o programa utiliza dados de 2 locais: [Binance](https://www.binance.com/pt-BR) 
para as smartcoins e [Investing.com](https://br.investing.com/) - para as ações.

    {  
            "ativos": [
		{
			"simbolo": "ETHBRL",
			"link": "https://api.binance.com/api/v3/ticker/24hr?symbol=ETHBRL",
			"rsi": "https://br.investing.com/crypto/ethereum/eth-brl-technical",
			"tipo": "criptomoeda",
			"taxa": 0.00,
			"quantidade": 0.03695740,
			"inicial": 22458.68,
			"perda": 0,
			"ganho": 0,
			"alerta_inf": 0,
			"alerta_sup": 0,
			"alerta_perc": 0
		},			
		{
			"simbolo": "IVVB11",
			"link": "https://br.investing.com/etfs/fundo-de-invest-ishares-sp-500",
			"rsi": "",
			"tipo": "etf",
			"quantidade": 58.00,
			"inicial": 234.22,
			"perda": 0,
			"ganho": 0,
			"alerta_perc": 0
		}	
	]
    }

**simbolo**: símbolo da moeda. Para obter uma listagem dos símbolos, [clique aqui](https://www.binance.com/api/v3/ticker/price)

**link**: é utilizado a URL *https://api.binance.com/api/v3/ticker/24hr?symbol=*
concatenada com o símbolo da smartcoin. Para saber mais: [API Binance](https://binance-docs.github.io/apidocs/spot/en/#introduction)

**rsi**: para extrair o valor do RSI, é utilizada a URL https://br.investing.com/crypto/ 
concatenada com o nome da smartcoin mais o complemento, que é formado pela sigla da 
smartcoin mais "-brl-technical". Exemplo: o Ethereum, é "ethereum/eth-brl-technical".

**tipo**: tipo do ativo: criptomoeda, acao ou etf

**taxa**: percentual a ser descontado do total, referente à alguma taxa que tenha que ser
para ao site no momento da negociação.

**quantidade**: quantidade de ativos, com ponto como separador decimal e sem separador de 
milhar.

**inicial**: valor inicial pago na aquisição do ativo (por unidade).

**perda**: valor limite para perda. Exemplo: 200.00. Ao atingir essa perda, é enviada uma 
mensagem.

**ganho**: valor de ganho desejado. Exemplo: 500.00. Ao atingir esse ganho, é enviada uma
mensagem.

**alerta_inf**: valor do limite inferior. Exemplo: 40000.00. Ao atingir esse valor, é 
enviada uma mensagem.  

**alerta_sup**: valor do limite superior. Exemplo: 60000.00. Ao atingir esse valor, é
enviada uma mensagem.

**alerta_perc**: percentual de ganho/perda. Exemplo: 1.00. Ao ter um ganho de 1%, é enviada 
uma mensagem. Para ativar a mensagem na perda, use percentuais negativos: -0.50

