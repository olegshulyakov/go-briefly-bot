package org.springframework.context.support;

import org.apache.commons.io.input.ReaderInputStream;
import org.springframework.beans.factory.config.YamlPropertiesFactoryBean;
import org.springframework.core.io.InputStreamResource;
import org.springframework.util.PropertiesPersister;

import java.io.*;
import java.util.Properties;

/**
 * A component responsible for loading properties from YAML files.
 * This class implements the {@link PropertiesPersister} interface,
 * providing functionality to load properties from YAML-formatted input streams and readers.
 * It explicitly does *not* support storing properties or handling XML formats, throwing
 * {@link UnsupportedOperationException} for such operations.
 */
public class YamlPropertiesLoader implements PropertiesPersister {

    /**
     * Loads properties from a YAML-formatted input stream.
     *
     * @param props The {@link Properties} object to populate with the loaded properties.
     * @param is   The {@link InputStream} containing the YAML data.
     */
    @Override
    public void load(Properties props, InputStream is) {
        YamlPropertiesFactoryBean yaml = new YamlPropertiesFactoryBean();
        yaml.setResources(new InputStreamResource(is));
        props.putAll(yaml.getObject());
    }

    /**
     * Loads properties from a YAML-formatted reader.
     *
     * @param props The {@link Properties} object to populate with the loaded properties.
     * @param reader The {@link Reader} containing the YAML data.
     * @throws IOException if an error occurs during reader reading or YAML parsing.
     */
    @Override
    public void load(Properties props, Reader reader) throws IOException {
        // Uses Commons IO ReaderInputStream
        InputStream inputStream = ReaderInputStream.builder().setReader(reader).get();
        load(props, inputStream);
    }

    /**
     * Throws an UnsupportedOperationException to indicate that storing properties is not supported.
     *
     * @param props The Properties object. (Not used)
     * @param os   The OutputStream. (Not used)
     * @param header The header. (Not used)
     * @throws UnsupportedOperationException Always thrown.
     * @throws IOException  If an IOException occurs during OutputStream operation.
     */
    @Override
    public void store(Properties props, OutputStream os, String header) throws IOException {
        throw new UnsupportedOperationException("Storing properties is not supported by YamlPropertiesLoader");
    }

    /**
     * Throws an UnsupportedOperationException to indicate that storing properties is not supported.
     *
     * @param props The Properties object. (Not used)
     * @param writer The Writer. (Not used)
     * @param header The header. (Not used)
     * @throws UnsupportedOperationException Always thrown.
     * @throws IOException If an IOException occurs during Writer operation
     */
    @Override
    public void store(Properties props, Writer writer, String header) throws IOException {
        throw new UnsupportedOperationException("Storing properties is not supported by YamlPropertiesLoader");
    }

    /**
     * Throws an UnsupportedOperationException to indicate that loading from XML is not supported.
     *
     * @param props The Properties object. (Not used)
     * @param is The InputStream. (Not used)
     * @throws UnsupportedOperationException Always thrown.
     * @throws IOException If an IOException occurs during InputStream operation
     */
    @Override
    public void loadFromXml(Properties props, InputStream is) throws IOException {
        throw new UnsupportedOperationException("Loading from XML is not supported by YamlPropertiesLoader");
    }

    /**
     * Throws an UnsupportedOperationException to indicate that storing to XML is not supported.
     *
     * @param props The Properties object. (Not used)
     * @param os The OutputStream. (Not used)
     * @param header The header. (Not used)
     * @throws UnsupportedOperationException Always thrown.
     * @throws IOException If an IOException occurs during OutputStream operation
     */
    @Override
    public void storeToXml(Properties props, OutputStream os, String header) throws IOException {
        throw new UnsupportedOperationException("Storing to XML is not supported by YamlPropertiesLoader");
    }

    /**
     * Throws an UnsupportedOperationException to indicate that storing to XML is not supported.
     *
     * @param props The Properties object. (Not used)
     * @param os The OutputStream. (Not used)
     * @param header The header. (Not used)
     * @param encoding The encoding. (Not used)
     * @throws UnsupportedOperationException Always thrown.
     * @throws IOException If an IOException occurs during OutputStream operation
     */
    @Override
    public void storeToXml(Properties props, OutputStream os, String header, String encoding) throws IOException {
        throw new UnsupportedOperationException("Storing to XML is not supported by YamlPropertiesLoader");
    }
}
