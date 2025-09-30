# PW Equipment Changer

Um aplicativo para trocar automaticamente equipamentos no Perfect World usando teclas de atalho.

## Funcionalidades

- Interface gráfica moderna usando Fyne
- Interface de linha de comando (console)
- **Sistema de assinatura SAAS** com verificação por HWID
- Configuração de até 11 itens para troca
- Tempo configurável entre cliques
- Monitoramento de tecla Q para ativar a troca
- Verificação automática de assinatura antes do uso

## Como usar

### Interface Gráfica (Recomendado)

```bash
./pw-equip-gui
```

ou

```bash
./pw-equip-gui -gui=true
```

### Interface de Console

```bash
./pw-equip-gui -gui=false
```

### Mostrar HWID da Máquina

```bash
./pw-equip-gui -hwid
```

**Nota**: O HWID é necessário para ativar a assinatura. Entre em contato com o suporte fornecendo este código.

## Configuração

### Campos necessários:

1. **Número de itens**: Quantos equipamentos você deseja trocar (1-11)
2. **Tecla para mudar barras**: A tecla que muda as barras de skills (`v` ou `` ` ``)
3. **Tempo entre cliques**: Tempo em milissegundos entre cada clique nos itens
4. **Teclas dos itens**: As teclas correspondentes a cada item que será trocado

### Instruções de setup:

1. Deixe 3 barras livres para serem rotacionadas
2. Em sua barra principal, deixe suas skills/boticários como deseja usá-los
3. Se deseja iniciar com equipamentos de ataque, na segunda barra deixe os Equipamentos de ataque
4. Na última barra, deixe os Equipamentos de defesa
5. Para trocar de set, aperte a tecla **Q**!

## Compilação

### Compilação Simples
```bash
go build -o pw-equip-gui
```

### Compilação Multiplataforma
```bash
# Use o script de build automático
./build.sh
```

### Para macOS - Evitar Janela do Terminal
Se o Terminal abrir quando você clicar no executável no macOS:

```bash
# 1. Compile normalmente
go build -o pw-equip-gui

# 2. Crie um App Bundle
./create_app_bundle.sh

# 3. Use o arquivo .app gerado
```

### Para Windows - Evitar Janela do Console
```bash
# No Windows, use a flag especial
go build -ldflags="-H windowsgui" -o pw-equip-gui.exe
```

## Sistema de Assinatura

O aplicativo requer uma assinatura ativa para funcionar. O sistema funciona da seguinte forma:

1. **HWID (Hardware ID)**: Cada máquina possui um identificador único baseado em informações do hardware
2. **Verificação**: Antes de iniciar o monitoramento, o app verifica se a assinatura está ativa
3. **API**: Faz uma requisição GET para `http://200.1.1.1/subscription?hwid=<HWID>`
4. **Resposta**: A API retorna `{"status": true/false}` indicando se a assinatura está ativa

### Para ativar sua assinatura:

1. Execute `./pw-equip-gui -hwid` para obter seu HWID
2. Entre em contato com o suporte fornecendo o HWID
3. Após a ativação, o aplicativo funcionará normalmente

## Dependências

- [Fyne](https://fyne.io/) - Interface gráfica
- [robotgo](https://github.com/go-vgo/robotgo) - Automação de teclado
- [gohook](https://github.com/robotn/gohook) - Captura de teclas
- [gopsutil](https://github.com/shirou/gopsutil) - Informações do sistema para HWID

## Bibliotecas GUI Consideradas

Durante o desenvolvimento, foram avaliadas as seguintes opções:

1. **[Fyne](https://github.com/fyne-io/fyne)** ✅ **Escolhida**
   - Fácil de usar e aprender
   - Excelente documentação
   - Deploy em binário único
   - Perfeita para formulários simples

2. **[GoVCL](https://github.com/ying32/govcl)**
   - Interface nativa
   - Requer arquivos DLL/SO externos
   - Mais complexa para setup

3. **[Spot](https://github.com/roblillack/spot)**
   - Abordagem moderna similar ao React
   - Menos madura, comunidade menor
