/**
 * Product Controller
 */
const Product = require('../models/Product');
const Category = require('../models/Category');
const { Op } = require('sequelize');

// Helper function to generate slug
const generateSlug = (name) => {
  return name
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/(^-|-$)/g, '');
};

/**
 * Get all products with filtering, search, and pagination
 */
exports.getAllProducts = async (req, res) => {
  try {
    const {
      page = 1,
      limit = 20,
      search = '',
      category_id,
      min_price,
      max_price,
      is_featured,
      sort_by = 'created_at',
      order = 'DESC'
    } = req.query;

    // Build where clause
    const where = { is_active: true };

    // Search by name or description
    if (search) {
      where[Op.or] = [
        { name: { [Op.iLike]: `%${search}%` } },
        { description: { [Op.iLike]: `%${search}%` } },
        { tags: { [Op.contains]: [search] } }
      ];
    }

    // Filter by category
    if (category_id) {
      where.category_id = category_id;
    }

    // Filter by price range
    if (min_price || max_price) {
      where.price = {};
      if (min_price) where.price[Op.gte] = min_price;
      if (max_price) where.price[Op.lte] = max_price;
    }

    // Filter by featured
    if (is_featured !== undefined) {
      where.is_featured = is_featured === 'true';
    }

    // Calculate pagination
    const offset = (page - 1) * limit;

    // Query products
    const { count, rows } = await Product.findAndCountAll({
      where,
      include: [
        {
          model: Category,
          as: 'category',
          attributes: ['id', 'name', 'slug']
        }
      ],
      limit: parseInt(limit),
      offset: parseInt(offset),
      order: [[sort_by, order.toUpperCase()]],
      distinct: true
    });

    res.json({
      products: rows,
      pagination: {
        total: count,
        page: parseInt(page),
        limit: parseInt(limit),
        pages: Math.ceil(count / limit)
      }
    });

  } catch (error) {
    console.error('Error fetching products:', error);
    res.status(500).json({ error: 'Failed to fetch products' });
  }
};

/**
 * Get single product by ID or slug
 */
exports.getProduct = async (req, res) => {
  try {
    const { id } = req.params;

    // Check if id is a number or slug
    const where = isNaN(id) ? { slug: id } : { id: parseInt(id) };

    const product = await Product.findOne({
      where: { ...where, is_active: true },
      include: [
        {
          model: Category,
          as: 'category',
          attributes: ['id', 'name', 'slug']
        }
      ]
    });

    if (!product) {
      return res.status(404).json({ error: 'Product not found' });
    }

    res.json({ product });

  } catch (error) {
    console.error('Error fetching product:', error);
    res.status(500).json({ error: 'Failed to fetch product' });
  }
};

/**
 * Create new product
 */
exports.createProduct = async (req, res) => {
  try {
    const {
      name,
      description,
      short_description,
      price,
      compare_at_price,
      cost_per_item,
      sku,
      barcode,
      quantity,
      category_id,
      images,
      weight,
      dimensions,
      tags,
      is_featured,
      meta_title,
      meta_description
    } = req.body;

    // Generate slug from name
    const slug = generateSlug(name);

    // Check if slug already exists
    const existingProduct = await Product.findOne({ where: { slug } });
    if (existingProduct) {
      return res.status(400).json({ error: 'Product with this name already exists' });
    }

    // Create product
    const product = await Product.create({
      name,
      slug,
      description,
      short_description,
      price,
      compare_at_price,
      cost_per_item,
      sku,
      barcode,
      quantity: quantity || 0,
      category_id,
      images: images || [],
      weight,
      dimensions,
      tags: tags || [],
      is_featured: is_featured || false,
      meta_title,
      meta_description
    });

    // Fetch with category
    const productWithCategory = await Product.findByPk(product.id, {
      include: [
        {
          model: Category,
          as: 'category',
          attributes: ['id', 'name', 'slug']
        }
      ]
    });

    res.status(201).json({
      message: 'Product created successfully',
      product: productWithCategory
    });

  } catch (error) {
    console.error('Error creating product:', error);
    res.status(500).json({ error: 'Failed to create product' });
  }
};

/**
 * Update product
 */
exports.updateProduct = async (req, res) => {
  try {
    const { id } = req.params;
    const updateData = req.body;

    // Find product
    const product = await Product.findByPk(id);
    if (!product) {
      return res.status(404).json({ error: 'Product not found' });
    }

    // Update slug if name changed
    if (updateData.name && updateData.name !== product.name) {
      updateData.slug = generateSlug(updateData.name);
    }

    // Update product
    await product.update(updateData);

    // Fetch updated product with category
    const updatedProduct = await Product.findByPk(id, {
      include: [
        {
          model: Category,
          as: 'category',
          attributes: ['id', 'name', 'slug']
        }
      ]
    });

    res.json({
      message: 'Product updated successfully',
      product: updatedProduct
    });

  } catch (error) {
    console.error('Error updating product:', error);
    res.status(500).json({ error: 'Failed to update product' });
  }
};

/**
 * Delete product (soft delete)
 */
exports.deleteProduct = async (req, res) => {
  try {
    const { id } = req.params;

    const product = await Product.findByPk(id);
    if (!product) {
      return res.status(404).json({ error: 'Product not found' });
    }

    // Soft delete by setting is_active to false
    await product.update({ is_active: false });

    res.json({
      message: 'Product deleted successfully'
    });

  } catch (error) {
    console.error('Error deleting product:', error);
    res.status(500).json({ error: 'Failed to delete product' });
  }
};

/**
 * Get featured products
 */
exports.getFeaturedProducts = async (req, res) => {
  try {
    const { limit = 10 } = req.query;

    const products = await Product.findAll({
      where: {
        is_active: true,
        is_featured: true
      },
      include: [
        {
          model: Category,
          as: 'category',
          attributes: ['id', 'name', 'slug']
        }
      ],
      limit: parseInt(limit),
      order: [['created_at', 'DESC']]
    });

    res.json({ products });

  } catch (error) {
    console.error('Error fetching featured products:', error);
    res.status(500).json({ error: 'Failed to fetch featured products' });
  }
};