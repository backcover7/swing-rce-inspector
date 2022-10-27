# swing-rce-inspector

自用小工具，没有什么太大的意义，练习`Golang`为主

本项目编写于`2022/9/27`，在`2022/10/7`日上传Github，于`2022/10/28`开源

## 对比分析

网上已有一些分析工具，本项目的优势在于使用简单效率较高，可拓展性远不如其他工具

同时分析`rt.jar`与`cobaltstrike.jar`文件速度对比：

- 使用tabby分析：初始化时间很长
- 使用CodeQL分析：初始化时间很长
- 使用Java ASM分析：在`Linux`下`6-8`秒，在`Windows`下`18-20`秒
- 使用本项目分析：在`Linux`下`2-3`秒，在`Windows`下`7-9`秒

## 使用

纯`Golang`编写，无任何其他依赖，编译生成可执行文件，快速分析

命令：`./swing-rce-inspector`

前提：

- 在`swing-rce-inspector`当前目录创建`jars`目录
- 请把需要的`Jar`文件放入`jars`目录中（例如`rt.jar`）
- 请以`root`权限运行

扫描`rt.jar`和`cobaltstrike.jar`的结果如下：

```text
符合条件的类 目标set方法
-> 该set方法中存在的方法调用1 (最有可能的方法会标记taint)
-> 该set方法中存在的方法调用2 (最有可能的方法会标记taint)
-> 该set方法中存在的方法调用3 (最有可能的方法会标记taint)
-> 该set方法中存在的方法调用n (最有可能的方法会标记taint)

org/apache/batik/apps/svgbrowser/StatusBar setMessage
-> org/apache/batik/apps/svgbrowser/StatusBar.getPreferredSize
-> java/awt/Dimension.<init>
-> org/apache/batik/apps/svgbrowser/StatusBar.setPreferredSize
-> org/apache/batik/apps/svgbrowser/StatusBar$DisplayThread.finish
-> org/apache/batik/apps/svgbrowser/StatusBar$DisplayThread.<init> (taint)
-> org/apache/batik/apps/svgbrowser/StatusBar$DisplayThread.start

org/apache/batik/apps/svgbrowser/StatusBar setMainMessage
-> javax/swing/JLabel.setText (taint)
-> org/apache/batik/apps/svgbrowser/StatusBar$DisplayThread.finish
-> org/apache/batik/apps/svgbrowser/StatusBar.getPreferredSize
-> java/awt/Dimension.<init>
-> org/apache/batik/apps/svgbrowser/StatusBar.setPreferredSize

org/apache/batik/swing/JSVGCanvas setURI
-> org/apache/batik/swing/JSVGCanvas.loadSVGDocument (taint)
-> org/apache/batik/swing/JSVGCanvas.setSVGDocument
-> java/beans/PropertyChangeSupport.firePropertyChange

org/apache/batik/swing/svg/AbstractJSVGComponent setFragmentIdentifier
-> org/apache/batik/swing/svg/AbstractJSVGComponent.computeRenderingTransform (taint)
-> org/apache/batik/swing/svg/AbstractJSVGComponent.scheduleGVTRendering

java/awt/Checkbox setLabel
-> java/lang/String.equals (taint)
-> java/awt/peer/CheckboxPeer.setLabel (taint)
-> java/awt/Checkbox.invalidateIfValid

java/awt/Label setText
-> java/lang/String.equals (taint)
-> java/awt/peer/LabelPeer.setText (taint)
-> java/awt/Label.invalidateIfValid

java/awt/Button setLabel
-> java/lang/String.equals (taint)
-> java/awt/peer/ButtonPeer.setLabel (taint)
-> java/awt/Button.invalidateIfValid

java/awt/Frame setTitle
-> java/awt/peer/FramePeer.setTitle (taint)
-> java/awt/Frame.firePropertyChange (taint)

javax/swing/AbstractButton setText
-> javax/swing/AbstractButton.firePropertyChange (taint)
-> javax/swing/AbstractButton.getMnemonic (taint)
-> javax/swing/AbstractButton.updateDisplayedMnemonicIndex
-> javax/accessibility/AccessibleContext.firePropertyChange (taint)
-> java/lang/String.equals (taint)
-> javax/swing/AbstractButton.revalidate
-> javax/swing/AbstractButton.repaint

javax/swing/JFileChooser setDialogTitle
-> javax/swing/JDialog.setTitle (taint)
-> javax/swing/JFileChooser.firePropertyChange (taint)

javax/swing/JFileChooser setApproveButtonToolTipText
-> javax/swing/JFileChooser.firePropertyChange (taint)

javax/swing/JFileChooser setApproveButtonText
-> javax/swing/JFileChooser.firePropertyChange (taint)

javax/swing/JInternalFrame setTitle
-> javax/swing/JInternalFrame.firePropertyChange (taint)

javax/swing/JLabel setText
-> javax/accessibility/AccessibleContext.getAccessibleName
-> javax/swing/JLabel.firePropertyChange (taint)
-> javax/swing/JLabel.getDisplayedMnemonic (taint)
-> javax/swing/SwingUtilities.findDisplayedMnemonicIndex
-> javax/swing/JLabel.setDisplayedMnemonicIndex
-> javax/accessibility/AccessibleContext.getAccessibleName
-> javax/accessibility/AccessibleContext.getAccessibleName
-> javax/accessibility/AccessibleContext.firePropertyChange
-> java/lang/String.equals (taint)
-> javax/swing/JLabel.revalidate
-> javax/swing/JLabel.repaint

javax/swing/JPopupMenu setLabel
-> javax/swing/JPopupMenu.firePropertyChange (taint)
-> javax/accessibility/AccessibleContext.firePropertyChange (taint)
-> javax/swing/JPopupMenu.invalidate
-> javax/swing/JPopupMenu.repaint

javax/swing/JToolTip setTipText
-> javax/swing/JToolTip.firePropertyChange (taint)
-> java/util/Objects.equals (taint)
-> javax/swing/JToolTip.revalidate
-> javax/swing/JToolTip.repaint
```

分析过程：

- 解压Jar包得到所有class文件
- 分析所有class文件得到类和方法信息
- 构建类之间的继承关系
- 根据规则分析得到结果

## 注意事项

如果想扫描多个`Jar`包，全部放入`jar`目录即可

**如果想扫描`cobaltstrike.jar`等第三方`jar`包，请保证`rt.jar`也加入了`lib`**

原因：

进行继承关系分析的时候，第三方`Jar`通常不包含`java.awt.Component`
类以及常用的子类，但实际上很多类是继承自它们的，如果不导入`rt.jar`会导致无法正确分析继承关系，以至于后续分析无法继续最终没有结果

继承关系分析原则：

```text
example:
class A extends B
class B extends C
class C implements D,E

result:
A is subclass of B,C,D,E
B is subclass of C,D,E
C is subclass of D,E
```
