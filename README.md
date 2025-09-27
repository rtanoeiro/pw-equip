# PW Equipment Changer

Um aplicativo para trocar automaticamente equipamentos no Perfect World usando teclas de atalho.

## Funcionalidades

- Interface gráfica moderna usando Fyne
- Interface de linha de comando (console)
- Configuração de até 11 itens para troca
- Tempo configurável entre cliques
- Monitoramento de tecla Q para ativar a troca

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

```bash
go build -o pw-equip-gui
```

## Dependências

- [Fyne](https://fyne.io/) - Interface gráfica
- [robotgo](https://github.com/go-vgo/robotgo) - Automação de teclado
- [gohook](https://github.com/robotn/gohook) - Captura de teclas

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
